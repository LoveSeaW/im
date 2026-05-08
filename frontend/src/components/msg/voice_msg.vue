<script setup lang="ts">
import type {voiceMsg} from "@/types/msg";
import {ref} from "vue";

interface Props {
  data: voiceMsg
  isMe: boolean
}

const props = defineProps<Props>()
const audio = ref<HTMLAudioElement>()

function play() {
  if (!audio.value) {
    // 刚开始
    audio.value = new Audio(props.data.src)
    audio.value?.play()
  } else {
    if (audio.value?.paused) {
      // 如果是暂停，就要播放
      audio.value?.play()
    } else {
      // 暂停
      audio.value?.pause()
    }

  }
}

</script>

<template>
  <div class="voice_msg" :class="{isMe: props.isMe}">
    <span>{{ props.data.time }}"</span>
    <i @click="play" class="iconfont icon-yuyinxiaoxi2x"></i>
  </div>
</template>

<style lang="scss">
.voice_msg {
  height: 40px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  color: var(--text_color);
  background-color: var(--item_hover);

  i {
    cursor: pointer;
  }

  &.isMe {
    background-color: var(--msg_active_color);
    color: white;
  }
}

</style>