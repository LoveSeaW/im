<script setup lang="ts">
import {ref} from "vue";
import {ElMessage} from "element-plus";
import {formatSecond} from "@/utils/file";
import {uploadFileApi} from "@/api/file_api";
import {useStore} from "@/stores";
import {Msg} from "@/service/send_msg";

const store = useStore()

interface Props {
  type: "user" | "group"
  id: number // 用户就是用户id，群就是群id
}

const props = defineProps<Props>()
const msg = new Msg(props.type)
const mediaRecorder = ref<MediaRecorder>()
const visible = ref(false)
const second = ref(0)
const timer = ref(0)
const isCancel = ref(false)


function voiceStart() {
  navigator.mediaDevices.getUserMedia({audio: true})
      .then(function (stream) {
        mediaRecorder.value = new MediaRecorder(stream);
        mediaRecorder.value.start()
        visible.value = true

        timer.value = setInterval(() => second.value++, 1000)

        // 录音结束时，停止并下载录音文件  stop方法会触发这个事件
        mediaRecorder.value.addEventListener("dataavailable", async function (event) {
          clearInterval(timer.value)
          if (isCancel.value) {
            ElMessage.success("取消录音")
            return
          }
          const audioFile = new Blob([event.data], {type: "audio/wav"});

          const res = await uploadFileApi(audioFile as File)
          if (res.code) {
            ElMessage.error(res.msg)
            return
          }
          // console.log(res.data)
          sendVoice(res.data.src)
          // console.log(audioFile)
          // const audioURL = URL.createObjectURL(audioFile);
          // const downloadLink = document.createElement("a");
          // downloadLink.href = audioURL;
          // downloadLink.download = "recording.mp3";
          // downloadLink.click();
        });
      })
      .catch(function (error) {
        // 用户不给权限
        ElMessage.error("用户取消录音权限")
        console.error(error);
      });
}

function sendVoice(src: string) {
  msg.voiceMsg(props.id, {
    src: src,
    time: second.value,
  })
  second.value = 0
}


function cancel() {
  // 结束录音
  // 时间归零
  second.value = 0
  isCancel.value = true
  mediaRecorder.value?.stop()
  visible.value = false
}

function send() {
  mediaRecorder.value?.stop()
  visible.value = false
}

function close() {

}


</script>

<template>
  <div class="fim_voice">
    <el-dialog modal-class="voice_dialog" @close-auto-focus="close" v-model="visible" width="40%">
      <div class="voice_body">
        <div class="progress">
          <el-progress :percentage="second">{{ formatSecond(second) }}</el-progress>
        </div>
        <div class="action">
          <el-button type="primary" size="small" @click="send">发送</el-button>
          <el-button size="small" @click="cancel">取消</el-button>
        </div>
      </div>
    </el-dialog>
    <i class="iconfont icon-luyin" title="语音消息" @click="voiceStart"></i>
  </div>
</template>

<style lang="scss">
.voice_dialog {
  .el-dialog__header {
    display: none;
  }

  .voice_body {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .progress {
      width: 70%;
    }
  }
}
</style>