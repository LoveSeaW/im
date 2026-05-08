<script setup lang="ts">
import {ElAvatar} from "element-plus";

interface Props {
  src: string
  size?: number | "large" | "default" | "small"
  shape?: "circle" | "square" // 图片的模式
  online?: boolean
}


const props = defineProps<Props>()
const {size = 40} = props
</script>

<template>
  <span class="fim_avatar" :style="{width: size + 'px', height: size + 'px'}">
     <ElAvatar class="fim_avatar"
                :shape="props.shape"
                v-if="props.src.startsWith('http') || props.src.startsWith('/api')"
                :src="props.src"
                :size="size"></ElAvatar>
  <ElAvatar class="fim_avatar"
             :shape="props.shape"
             v-else :size="size" style="background-color: #529b2e">{{
      props.src
    }}
  </ElAvatar>

    <span class="fim_online" v-if="props.online"></span>
  </span>
</template>

<style>
.fim_avatar {
  flex-shrink: 0;
  position: relative;

  .fim_online {
    position: absolute;
    width: 5px;
    height: 5px;
    display: block;
    background-color: #54e009;
    border-radius: 50%;
    right: 1px;
    bottom: 4px;
  }
}
</style>