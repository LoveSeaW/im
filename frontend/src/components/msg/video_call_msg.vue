<script setup lang="ts">
import {ElImage} from "element-plus";
import type {videoCallMsg} from "@/types/msg";
import {time_different} from "@/utils/date";

interface Props {
  data: videoCallMsg
  isMe?: boolean
}

const props = defineProps<Props>()
</script>

<template>
  <div class="video_call_msg" :class="{isMe: props.isMe}">
    <i class="iconfont icon-shipin"></i>
    <span class="label">视频通话结束</span>
    <span>{{ time_different(props.data.startTime as string, props.data.endTime as string) }}</span>
  </div>
</template>

<style lang="scss">
.video_call_msg {
  height: 40px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  padding: 0 10px;
  color: var(--text_color);
  background-color: var(--item_hover);
  .label{
    margin: 0 5px;
  }

  &.isMe {
    background-color: var(--msg_active_color);
    color: white;
  }
}
</style>