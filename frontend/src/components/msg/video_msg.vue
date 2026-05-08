<script setup lang="ts">
import {VideoPlayer, type VideoPlayerProps} from '@videojs-player/vue'
import 'video.js/dist/video-js.css'
import type {videoMsg} from "@/types/msg";
interface Props {
  data: videoMsg
}
const props = defineProps<Props>()
const options: VideoPlayerProps = {
  autoplay: false,
  loop: true,
  controls: true,
  disablePictureInPicture: false,
  playbackRates: [0.7, 1.0, 1.5, 2.0], // 播放速度
  width: 350,
  sources: [
    {
      type: 'video/mp4', // 资源格式写法：'video/mp4'
      src: props.data.src
    }
  ],
}
</script>

<template>
  <div class="video_msg">
    <video-player :options="options"/>
  </div>
</template>

<style scoped lang="scss">

</style>