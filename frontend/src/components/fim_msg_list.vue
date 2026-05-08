<script setup lang="ts">
import Fim_msg from "@/components/fim_msg.vue";
import type {msgChatType} from "@/types/msg";
import {ref} from "vue";
import {useRoute} from "vue-router";
import type {baseResponse} from "@/api";
import {chatDeleteApi} from "@/api/chat_api";
import {groupChatDeleteApi} from "@/api/group_api";
import {ElMessage} from "element-plus";
import type {fimMsgListExpose} from "@/types/msg_list";

const route = useRoute()

interface Props {
  msgList: msgChatType[]
  type: "user" | "group"
  showNickname?: boolean
}

const props = defineProps<Props>()

const emits = defineEmits(["checkShow"])

const checkShowVisible = ref(false)

function checkShow() {
  checkShowVisible.value = true
  emits("checkShow")
}

const useMsgIDList = ref<number[]>([])

function check(type: string, msgID: number) {
  if (type === "add") {
    useMsgIDList.value.push(msgID)
    return
  }
  if (type === "remove") {
    const index = useMsgIDList.value.findIndex((value) => value === msgID)
    if (index === -1) return;
    useMsgIDList.value.splice(index, 1)
  }
}

interface checkType {
  closeCheck: () => void
}

const fimMsgRef = ref<checkType[]>([])

function close() {
  useMsgIDList.value = []
  checkShowVisible.value = false
  for (const argument of fimMsgRef.value) {
    argument.closeCheck()
  }
}

async function chatDelete() {
  let res: baseResponse<string> = {code: 0, msg: "", data: ""}
  if (props.type === "user") {
    res = await chatDeleteApi(useMsgIDList.value)
  } else {
    res = await groupChatDeleteApi(Number(route.params.id), useMsgIDList.value)
  }
  if (res.code) {
    ElMessage.error(res.msg)
    return false
  }
  useMsgIDList.value = []
  if (props.type === "user") {
    ElMessage.success("删除对话聊天记录成功")
  } else {
    ElMessage.success("删除群聊天记录成功")
  }
  return true
}

defineExpose({
  close,
  chatDelete
} as fimMsgListExpose)


</script>

<template>
  <fim_msg :key="item.id" ref="fimMsgRef" @check="check" @check-show="checkShow"
           :check-show="checkShowVisible"
           :type="props.type" :show-nickname="props.showNickname" v-for="item in props.msgList" :data="item"></fim_msg>
</template>

<style scoped lang="scss">

</style>