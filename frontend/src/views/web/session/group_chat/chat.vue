<script setup lang="ts">
import type {msgChatType} from "@/types/msg";
import {reactive, type Ref, ref, watch} from "vue";
import type {
  groupHistoryListParams,
  groupHistoryType,
} from "@/api/group_api";
import {useRoute} from "vue-router";
import {groupHistoryListApi} from "@/api/group_api";
import {useStore} from "@/stores";
import {ElMessage} from "element-plus";
import Group_member from "@/components/group/group_member.vue";
import {Close, Delete} from "@element-plus/icons-vue";
import type {fimMsgListExpose} from "@/types/msg_list";
import Fim_msg_list from "@/components/fim_msg_list.vue";
import Fim_image_upload from "@/components/chat/fim_image_upload.vue";
import Fim_file_upload from "@/components/chat/fim_file_upload.vue";
import Fim_emoji from "@/components/chat/fim_emoji.vue";
import Fim_video_upload from "@/components/chat/fim_video_upload.vue";
import Fim_chat_ipt from "@/components/chat/fim_chat_ipt.vue";
import Fim_event from "@/components/fim_event.vue";
import Fim_voice from "@/components/chat/fim_voice.vue";

const store = useStore()
const route = useRoute()
const form = reactive({
  content: ""
})
const msgList = ref<msgChatType[]>([])


const params = reactive<groupHistoryListParams>({
  id: 0,
  limit: 10,
  page: 1,
})

async function getData() {
  params.id = Number(route.params.id)
  params.page = 1 // page复位
  let res = await groupHistoryListApi(params)
  const list = res.data.list || []
  msgList.value = []
  for (const re of list) {
    msgList.value.push(getMsg(re))
  }

}


function getMsg(msg: groupHistoryType): msgChatType {
  return {
    "id": msg.id,
    "user": {
      id: msg.userID,
      nickName: msg.userNickname,
      avatar: msg.userAvatar,
    },
    "isMe": msg.isMe,
    "created_at": msg.createdAt,
    "msg": msg.msg,
    msgPreview: msg.msgPreview,
    showDate: msg.showDate,
  }
}

const scrollbarRef = ref()
watch(() => store.groupMsgData, () => {
  if (store.groupMsgData.msg.type === 12) {
    ElMessage.error(store.groupMsgData.msg.tipMsg?.content)
    return
  }
  if (store.groupMsgData.groupID === Number(route.params.id)) {
    // 如果是撤回消息
    if (store.groupMsgData.msg.type === 8) {
      getData()
      return;
    }
    if (store.groupMsgData.id === 0) {
      return;
    }
    msgList.value.push(getMsg(store.groupMsgData))
    setTimeout(() => {
      const top = (document.querySelector(".group_msg_body") as HTMLDivElement).clientHeight
      scrollbarRef.value.setScrollTop(top)
    }, 100)
  }
}, {deep: true})


watch(() => route.params, () => {
  getData()
}, {immediate: true, deep: true})


const checkShowVisible = ref(false)
const fimMsgListRef = ref<fimMsgListExpose | undefined>()

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
  params.id = Number(route.params.id)
  // 如何判断点到头了

  params.page = (params.page as number) + 1
  let res = await groupHistoryListApi(params)
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
  <div class="group_chat_inner">
    <div class="left">
      <div class="group_chat_body">
        <el-scrollbar ref="scrollbarRef" height="100%">
          <div class="group_msg_body" @mousedown.middle.prevent="store.showMsgDateMethod">
            <div class="more" v-if="moreVisible">
              <span @click="more">加载更多</span>
            </div>
            <fim_msg_list ref="fimMsgListRef" @checkShow="checkShowVisible=true" show-nickname :msg-list="msgList"
                          type="group"></fim_msg_list>
          </div>
        </el-scrollbar>
      </div>
      <div class="group_chat_inner_action" v-if="checkShowVisible">
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
        <div class="group_msg_menu">
          <fim_emoji></fim_emoji>
          <fim_image_upload :id="Number(route.params.id)" type="group"></fim_image_upload>
          <fim_video_upload :id="Number(route.params.id)" type="group"></fim_video_upload>
          <fim_file_upload :id="Number(route.params.id)" type="group"></fim_file_upload>
          <fim_voice :id="Number(route.params.id)" type="group"></fim_voice>
        </div>
        <div class="group_msg_box">
          <div class="group_msg_tip" v-if="store.groupData.isProhibition && store.groupData.role === 3 || store.groupData.prohibitionTime != null">
            <span v-if="store.groupData.isProhibition && store.groupData.role === 3">当前群正在全员禁言中</span>
            <span v-else-if="store.groupData.prohibitionTime != null">当前用户正在被禁言</span>
          </div>
          <fim_chat_ipt :id="Number(route.params.id)" type="group"></fim_chat_ipt>
        </div>
      </template>


    </div>

    <fim_event event-key="history" @event="getData"></fim_event>
    <div class="right">

      <div class="group_notice">
        <div class="title">群通知</div>
        <div class="body">
          <el-empty description="暂无通知" :image-size="50"></el-empty>
        </div>
      </div>
      <group_member :group-id="Number(route.params.id)"></group_member>
    </div>
  </div>
</template>

<style lang="scss">
.group_chat_inner {
  width: 100%;
  height: 100%;
  display: flex;

  .right {
    width: 180px;

    .group_notice {
      width: 100%;
      height: 150px;
      padding: 10px;
      border-bottom: 1px solid var(--border_color);

      .title {
        font-size: 14px;
        color: var(--text_color);
      }

      .body {
        height: 100%;

        .el-empty {
          padding: 10px;
        }
      }
    }


  }

  .left {
    width: calc(100% - 180px);
    border-right: 1px solid var(--border_color);

    .group_chat_body {
      height: calc(100% - 150px);
      width: 100%;

      .more {
        display: flex;
        justify-content: center;
        font-size: 12px;
        color: var(--text_color);

        span {
          cursor: pointer;
        }
      }

      .el-scrollbar {
        width: 100%;
      }

      .group_msg_body {
        padding: 10px;
      }
    }

    .group_msg_menu {
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

    .group_msg_box {
      position: relative;

      .group_msg_tip {
        position: absolute;
        width: 100%;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 1;
      }

      .el-textarea__inner {
        box-shadow: none;
      }
    }

    .group_chat_inner_action {
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
}

.add_member_group_dialog {
  display: flex;
  justify-content: space-between;

  .dialog_left {
    width: 50%;
    border-right: 1px solid var(--border_color);
  }

  .dialog_right {
    width: 50%;
    padding-left: 10px;
  }

  .member_list {
    width: 100%;
    margin-top: 5px;

    .item {
      display: flex;
      align-items: center;

      .el-checkbox__label {
        display: flex;
        align-items: center;
      }

      .el-text {
        margin-left: 5px;
      }
    }

    .item_check {
      position: relative;
      height: 32px;

      i {
        opacity: 0;
        position: absolute;
        right: 0;
        cursor: pointer;
      }

      &:hover {
        i {
          opacity: 1;
        }
      }
    }
  }
}
</style>