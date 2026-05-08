<script setup lang="ts">
import {ref, watch} from "vue";
import {ElMessage} from "element-plus";
import {useStore} from "@/stores";
import {Msg} from "@/service/send_msg";

const visible = ref(false)
const localVideo = ref()
const localStream = ref<MediaStream>()
const store = useStore()
const remoteVideo = ref()
const status = ref(0) // 0  等待接听   1 对话中
interface Props {
  id: number // 用户就是用户id，群就是群id
  nickName: string // 用户的昵称
}
const msg = new Msg("user")

const props = defineProps<Props>()


// 发起视频通话 videoCall
function videoCall() {
  visible.value = true
  navigator.mediaDevices.getUserMedia({audio: true, video: true}).then((mediaStream) => {
    (localVideo.value as HTMLVideoElement).srcObject = mediaStream;
    // 发送视频通话消息
    localStream.value = mediaStream
    localVideo.value.muted = false // 解除静音

    msg.videoCallMsg(props.id, {
      flag: 0
    })

    pc.value = new RTCPeerConnection()

    // 回复answer_ice
    pc.value.addEventListener('icecandidate', (event: any) => {
      if (event.candidate) {
        msg.videoCallMsg(props.id, {
          flag: 9,
          type: "answer_ice",
          data: event.candidate,
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

watch(() => store.chatMsgData, () => {
  // 对方不在线
  if (!visible.value) return;
  if (store.chatMsgData.msg.type === 12) {
    setTimeout(() => {
      visible.value = false
    }, 1000)
    return
  }
  if (store.chatMsgData.msg.type === 7) {

    if (store.chatMsgData.msg.videoCallMsg?.flag === 6) {
      // 对方挂断了
      ElMessage.info("对方结束视频通话")
      visible.value = false
      status.value = 0
      remoteVideo.value.srcObject = null
      localVideo.value.srcObject = null
      return;
    }

    if (store.chatMsgData.msg.videoCallMsg?.flag === 4) {
      // 对方挂断了
      ElMessage.info("用户拒绝了你的视频通话")
      visible.value = false
      status.value = 0
      remoteVideo.value.srcObject = null
      localVideo.value.srcObject = null
      return;
    }

    switch (store.chatMsgData.msg.videoCallMsg?.type) {
      case "create_offer":
        sendOffer()
        break;
      case "offer_ice":
        pc.value.addIceCandidate(store.chatMsgData.msg.videoCallMsg?.data);
        break
      case "answer":
        recvAnswer(store.chatMsgData.msg.videoCallMsg?.data);
        break
    }
  }


}, {deep: true})


const pc = ref();


// 回复answer
function recvAnswer(answer: any) {
  pc.value.setRemoteDescription(answer);
}

// 发送offer
function sendOffer() {
  // 发送方发offer，对方要回offer_ice
  pc.value.addTrack(localStream.value?.getVideoTracks()[0], localStream.value);
  pc.value.createOffer({offerToReceiveAudio: true, offerToReceiveVideo: true}).then((offer: any) => {
    pc.value.setLocalDescription(offer).then(() => {
      msg.videoCallMsg(props.id, {
        flag: 9,
        type: "offer",
        data: offer,
      })
    })
  })
}

function stopVideoCall() {
  // 结束视频通话
  pc.value.close()
  msg.videoCallMsg(props.id, {
    flag: status.value === 1 ? 4 : 1
  })


  visible.value = false
  ElMessage.info("结束视频通话")
  status.value = 0
  remoteVideo.value.srcObject = null
  localVideo.value.srcObject = null
}

</script>

<template>
  <div>
    <el-dialog width="400px" append-to-body v-model="visible" modal-class="video_call_dialog" draggable>
      <div class="video_action">
        <div class="info" v-if="status === 0">
          <div class="name">{{ props.nickName }}</div>
          <div class="label">等待接听</div>
        </div>
        <div class="action">
          <section v-if="status === 0">
            <span class="icon"><i class="iconfont icon-shipin"></i></span>
            <span>摄像头已开</span>
          </section>
          <section v-if="status === 0">
            <span class="icon"><i class="iconfont icon-shipin"></i></span>
            <span>麦克风已开</span>
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
    <i @click="videoCall" title="视频通话" class="iconfont icon-020-videocallingapp"></i>
  </div>
</template>

<style lang="scss">
.video_call_dialog {
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