<script setup lang="ts">
import {reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import {watch} from "vue";
import {useStore} from "@/stores";
import {Msg} from "@/service/send_msg";

const store = useStore()
const visible = ref(false)
const localVideo = ref()
const remoteVideo = ref()
const localStream = ref<MediaStream>()
const sendUser = reactive({
  nickName: "",
  userID: 0,
})
const msg = new Msg("user")
const status = ref(0) // 0  等待接听   1 对话中


watch(() => store.chatMsgData, () => {
  if (store.chatMsgData.msg.type === 7) {
    if (store.chatMsgData.msg.videoCallMsg?.flag === 2) {
      sendUser.userID = store.chatMsgData.sendUser.id
      sendUser.nickName = store.chatMsgData.sendUser.nickName
      visible.value = true
    }

    if (!visible.value) return;

    if (store.chatMsgData.msg.videoCallMsg?.flag === 6) {
      // 对方挂断了
      ElMessage.info("对方结束视频通话")
      visible.value = false
      status.value = 0
      remoteVideo.value.srcObject = null
        localVideo.value.srcObject = null
      return;
    }
    if (store.chatMsgData.msg.videoCallMsg?.flag === 3) {
      // 对方挂断了
      ElMessage.info("发起者已挂断")
      visible.value = false
      status.value = 0
      remoteVideo.value.srcObject = null
        localVideo.value.srcObject = null
      return;
    }

    switch (store.chatMsgData.msg.videoCallMsg?.type) {
      case "offer":
        recvOffer(store.chatMsgData.msg.videoCallMsg?.data)
        break;
      case "answer_ice":
        pc.value.addIceCandidate(store.chatMsgData.msg.videoCallMsg?.data);
        break
    }

  }
}, {deep: true})

const pc = ref()


// 接收方接收offer
function recvOffer(offer: any) {
  pc.value.addTrack(localStream.value?.getVideoTracks()[0], localStream.value);
  pc.value.setRemoteDescription(offer).then(() => {
    pc.value.createAnswer().then((answer: any) => {
      pc.value.setLocalDescription(answer).then(() => {
        msg.videoCallMsg(sendUser.userID, {
          "flag": 9,
          "type": "answer",
          "data": answer,
        })
      })
    })
  })
}

function startOffer() {
  // 获取媒体信息
  navigator.mediaDevices.getUserMedia({audio: true, video: true}).then((mediaStream) => {
    (localVideo.value as HTMLVideoElement).srcObject = mediaStream;
    // 发送视频通话消息
    localStream.value = mediaStream
    localVideo.value.muted = false

    // 发送视频通话消息
    // 让发送者开始去发offer
    msg.videoCallMsg(sendUser.userID, {
      "flag": 3,
    })

    pc.value = new RTCPeerConnection();
// 回复offer
    pc.value.addEventListener('icecandidate', (event: any) => {
      if (event.candidate) {
        msg.videoCallMsg(sendUser.userID, {
          "flag": 9,
          "type": "offer_ice",
          "data": event.candidate,
        })
      }
    })

    pc.value.addEventListener("track", (event: any) => {
      remoteVideo.value.srcObject = event.streams[0];
      status.value = 1
      remoteVideo.value.muted = false
    })


  }).catch((err) => {
    ElMessage.error("无权限获取视频信息")
  })
}

function stopVideoCall() {
  // 结束视频通话
  pc.value?.close()
  msg.videoCallMsg(sendUser.userID, {
    "flag": status.value === 1 ? 4 : 2,
  })

  visible.value = false
  status.value = 0
  remoteVideo.value.srcObject = null
  localVideo.value.srcObject = null
  ElMessage.info("结束视频通话")
}

</script>

<template>
  <div>
    <el-dialog width="400px" append-to-body v-model="visible" modal-class="video_call_rev_dialog" draggable>
      <div class="video_action">
        <div class="info" v-if="status===0">
          <div class="name">来自{{ sendUser.nickName }}的视频通话</div>
          <div class="label">等待接听</div>
        </div>
        <div class="action">
          <section v-if="status===0" @click="startOffer">
            <span class="icon green"><i class="iconfont icon-shipin"></i></span>
            <span>接听</span>
          </section>
          <section @click="stopVideoCall">
            <span class="icon red"><i class="iconfont icon-shipin"></i></span>
            <span>挂断</span>
          </section>
        </div>
        <video class="local_video" ref="localVideo" autoplay muted></video>
        <video class="remote_video" ref="remoteVideo" autoplay muted></video>
      </div>
    </el-dialog>
  </div>
</template>

<style lang="scss">
.video_call_rev_dialog {
  .el-dialog {
    padding: 0;
  }

  .el-dialog__header {
    height: 400px;

    .el-dialog__headerbtn {
      display: none;
    }
  }

  .el-dialog__body {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    pointer-events: none;
    width: 100%;
    height: 100%;

    .video_action {
      width: 100%;
      height: 100%;
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      color: white;
      background-color: #c4c4c4;

      .action {
        pointer-events: auto;
        z-index: 2;
        width: 100%;
        display: flex;
        justify-content: space-evenly;
        position: absolute;
        bottom: 40px;

        section {
          display: flex;
          flex-direction: column;
          align-items: center;

          .icon {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            background-color: var(--item_hover);
            display: flex;
            justify-content: center;
            align-items: center;
            margin-bottom: 5px;

            i {
              font-size: 22px;
              font-weight: 600;
            }
          }

          .red {
            background-color: #d71345;
            color: white;
          }

          .green {
            background-color: #11c24f;
            color: white;
          }
        }
      }

      .info {
        position: relative;
        z-index: 2;
        text-align: center;

        .name {
          font-size: 20px;
        }
      }

      .local_video {
        z-index: 1;
        width: 100%;
        height: 100%;
        position: absolute;
        top: 0;
        object-fit: cover;
      }

      .remote_video {
        position: absolute;
        width: 100px;
        height: 150px;
        border-radius: 5px;
        top: 10px;
        right: 10px;
        object-fit: cover;
        z-index: 2;
      }
    }
  }
}
</style>