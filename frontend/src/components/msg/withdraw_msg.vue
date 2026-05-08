<script setup lang="ts">
import type {msgChatType} from "@/types/msg";

interface Props {
  data: msgChatType
  isMe?: boolean
}

const props = defineProps<Props>()
</script>

<template>
  <div class="withdraw_msg">
    <span v-if="props.isMe">
      你{{props.data.msg.withdrawMsg?.content}}
    </span>
    <span v-else>
      {{ props.data.user.nickName }} {{props.data.msg.withdrawMsg?.content}}
    </span>
  </div>
</template>

<style lang="scss">
.withdraw_msg {
  text-align: center;
  width: 100%;

  span {
    background-color: var(--item_hover);
    border-radius: 5px;
    font-size: 12px;
    padding: 5px 10px;
    color: var(--text_color);
  }
}
</style>