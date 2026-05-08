<script setup lang="ts">
import type {replyMsg} from "@/types/msg";
import {dateTimeFormat} from "../../utils/date";

interface Props {
  data: replyMsg
  isMe?: boolean
}

const props = defineProps<Props>()


function goTop() {
  const dom = document.querySelector(`[msg_id="${props.data.msgID}"]`)
  const msgDom = dom?.parentElement?.parentElement?.parentElement?.parentElement
  if (msgDom) {
    const scrollBar = document.querySelector(".chat_body_scrollbar")
    if (scrollBar) {
      scrollBar.scrollTop = msgDom.offsetTop
    }
  }

}

</script>

<template>
  <div class="reply_msg" :class="{isMe: props.isMe}">
    <div class="reply_origin_msg" :reply_msg_id="props.data.msgID">
      <div class="reply_msg_info">
        <span class="name">{{ props.data.userNickName }}</span>
        <span class="date">{{ dateTimeFormat(props.data.originMsgDate as string) }}</span>
        <i @click="goTop" class="iconfont icon-zhiding"></i>
      </div>
      <div class="reply_content">{{ props.data.replyMsgPreview }}</div>
    </div>
    <div class="content">
      {{ props.data.content }}
    </div>
  </div>
</template>

<style lang="scss">
.reply_msg {
  color: var(--text_color);
  background-color: var(--item_hover);
  padding: 5px;
  border-radius: 5px;

  .content {
    margin: 5px 0;
  }

  .reply_origin_msg {
    padding: 5px;
    border-radius: 2px;
    background-color: var(--item_hover_1);

    .reply_content {
      margin: 5px 0;
    }

    .reply_msg_info {
      span:nth-child(1) {
        margin-right: 5px;
      }
    }

    .icon-zhiding {
      cursor: pointer;
    }
  }

  &.isMe {
    background-color: var(--msg_active_color);
    color: white;

    .reply_origin_msg {
      background-color: #2fc061;
    }
  }
}
</style>