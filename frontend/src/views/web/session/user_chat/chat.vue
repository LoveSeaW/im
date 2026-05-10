<script setup lang="ts">
import {reactive, ref, watch} from "vue";
import {chatHistoryListApi, type chatHistoryListParams, type chatHistoryType} from "@/api/chat_api";
import {useRoute} from "vue-router";
import {ElMessage} from "element-plus";
import {useStore} from "@/stores";
import {Delete} from "@element-plus/icons-vue";
import {Close} from "@element-plus/icons-vue";
import {eventStore} from "@/stores/event_store";
import Fim_msg_list from "@/components/fim_msg_list.vue";
import type {msgChatType} from "@/types/msg";
import type {Ref} from "vue/dist/vue";
import type {fimMsgListExpose} from "@/types/msg_list";
import Fim_image_upload from "@/components/chat/fim_image_upload.vue";
import Fim_file_upload from "@/components/chat/fim_file_upload.vue";
import Fim_video_upload from "@/components/chat/fim_video_upload.vue";
import Fim_chat_ipt from "@/components/chat/fim_chat_ipt.vue";
import Fim_emoji from "@/components/chat/fim_emoji.vue";
import Fim_voice from "@/components/chat/fim_voice.vue";
import Fim_video_call from "@/components/chat/fim_video_call.vue";
import {friendDetailApi, type friendType} from "@/api/user_api";
import Fim_event from "@/components/fim_event.vue";

const store = useStore()
const es = eventStore()
const route = useRoute()
const form = reactive({
  content: ""
})

const msgList = ref<msgChatType[]>([])


const params = reactive<chatHistoryListParams>({
  friendId: 0,
  limit: 10,
  page: 1,
})

async function getData() {
  params.friendId = Number(route.params.id)
  params.page = 1 // page复位
  let res = await chatHistoryListApi(params)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  const list = res.data.list || []
  // list.reverse()

  msgList.value = []
  for (const re of list || []) {
    const item: msgChatType = getMsg(re)
    msgList.value.push(item)
  }

  setTimeout(() => {
    const top = (document.querySelector(".msg_body") as HTMLDivElement).clientHeight
    scrollbarRef.value.setScrollTop(top)
  }, 100)
}

const data = reactive<friendType>({
  userID: 0,
  nickname: "",
  abstract: "",
  avatar: "",
  notice: "",
  isOnline: false,
})

async function getFriendData() {
  let res = await friendDetailApi({friendID: Number(route.params.id)})
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  Object.assign(data, res.data)
}

watch(() => route.params, () => {
  getFriendData()
}, {immediate: true, deep: true})

function getMsg(msg: chatHistoryType): msgChatType {
  return {
    "id": msg.id,
    "user": msg.sendUser,
    "isMe": msg.isMe,
    "created_at": msg.created_at,
    "msg": msg.msg,
    msgPreview: msg.msgPreview,
    showDate: msg.showDate
  }
}


watch(() => route.params, () => {
  getData()
  store.replayChatMsgData = undefined

}, {immediate: true, deep: true})


watch(() => store.chatMsgData, () => {
  // 判断这里收到的消息，必须得是当下我在和这个人聊的消息
  if (store.chatMsgData.msg.type === 12) {
    ElMessage.error(store.chatMsgData.msg.tipMsg?.content)
    return
  }
  const friendID = Number(route.params.id)
  if (store.chatMsgData.revUser.id === friendID || store.chatMsgData.sendUser.id === friendID) {
    // 如果它没有id，或者类型是tipMsg
    // 判断是不是对方撤回了消息
    if (store.chatMsgData.msg.type === 8) {
      getData()
      return;
    }


    if (store.chatMsgData.id === 0) {
      return;
    }
    msgList.value.push(getMsg(store.chatMsgData))

    setTimeout(() => {
      const top = (document.querySelector(".msg_body") as HTMLDivElement).clientHeight
      scrollbarRef.value.setScrollTop(top)
    }, 100)
  }
}, {deep: true})

const scrollbarRef = ref()
const checkShowVisible = ref(false)
const fimMsgListRef = ref<fimMsgListExpose>()

function closeCheck() {
  fimMsgListRef.value?.close()
  checkShowVisible.value = false
}


async function chatDelete() {
  if (await fimMsgListRef.value?.chatDelete()) {
    getData()
  }
}

const moreVisible = ref(true)

async function more() {
  params.friendId = Number(route.params.id)
  // 如何判断点到头了

  params.page = (params.page as number) + 1
  let res = await chatHistoryListApi(params)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }

  if (msgList.value.length >= res.data.count) {
    console.log("点到头了")
    ElMessage.info("已经是最早的消息了")
    moreVisible.value = false
    return
  }

  const list = res.data.list || []
  const msgListData = []
  for (const re of list || []) {
    const item: msgChatType = getMsg(re)
    msgListData.push(item)
  }

  msgList.value = [...msgListData, ...msgList.value]

}


</script>

<template>
  <div class="user_chat_inner_view">
    <div class="user_chat_inner_head">
      <el-scrollbar wrap-class="chat_body_scrollbar" ref="scrollbarRef" height="100%">
        <div class="msg_body" @mousedown.middle.prevent="store.showMsgDateMethod">
          <div class="more" v-if="moreVisible">
            <span @click="more">加载更多</span>
          </div>
          <fim_msg_list ref="fimMsgListRef" @checkShow="checkShowVisible=true" :msg-list="msgList"
                        type="user"></fim_msg_list>
        </div>
      </el-scrollbar>
    </div>
    <fim_event event-key="history" @event="getData"></fim_event>
    <div class="user_chat_inner_action" v-if="checkShowVisible">
      <div class="icon">
        <el-icon @click="chatDelete" :size="20">
          <Delete/>
        </el-icon>
        <span>删除</span>
      </div>
      <div class="icon close_icon">
        <el-icon @click="closeCheck" :size="20">
          <Close/>
        </el-icon>
      </div>
    </div>
    <template v-else>
      <div class="user_chat_inner_menu">
        <fim_emoji></fim_emoji>
        <fim_image_upload :id="Number(route.params.id)" type="user"></fim_image_upload>
        <fim_video_upload :id="Number(route.params.id)" type="user"></fim_video_upload>
        <fim_file_upload :id="Number(route.params.id)" type="user"></fim_file_upload>
        <fim_voice :id="Number(route.params.id)" type="user"></fim_voice>
        <fim_video_call :id="Number(route.params.id)" :nick-name="data.nickname"></fim_video_call>
      </div>
      <div class="user_chat_inner_box">
        <fim_chat_ipt :id="Number(route.params.id)" type="user"></fim_chat_ipt>
      </div>
    </template>

  </div>
</template>

<style lang="scss">
.user_chat_inner_view {
  height: 100%;

  .user_chat_inner_head {
    height: calc(100% - 190px);

    .msg_body {
      padding: 10px;

      .more {
        display: flex;
        justify-content: center;
        font-size: 12px;
        color: var(--text_color);

        span {
          cursor: pointer;
        }
      }
    }
  }

  .user_chat_inner_menu {
    height: 30px;
    display: flex;
    border-top: 1px solid var(--border_color);

    i {
      width: 30px;
      height: 30px;
      font-size: 18px;
      display: flex;
      justify-content: center;
      align-items: center;
      cursor: pointer;

      &:hover {
        background-color: var(--item_hover);
      }
    }
  }

  .user_chat_inner_box {
    .el-textarea__inner {
      box-shadow: none;
    }
  }

  .user_chat_inner_action {
    height: 110px;
    border-top: 1px solid var(--border_color);
    padding: 10px;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;

    .icon {
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;

      .el-icon {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        background-color: var(--item_hover);
      }

      span {
        font-size: 14px;
        color: var(--text_color);
      }
    }

    .close_icon {
      position: absolute;
      right: 20px;
    }

    i {
      cursor: pointer;
    }
  }
}

</style>