<script setup lang="ts">

import Fim_slide_head from "@/components/fim_slide_head.vue";
import {computed, reactive} from "vue";
import type {listResponse} from "@/api";
import Fim_avatar from "@/components/fim_avatar.vue";
import {relativeCurrentTime} from "@/utils/date";
import {chatSessionListApi} from "@/api/chat_api";
import {groupSessionListApi} from "@/api/group_api";
import router from "@/router";
import {useRoute} from "vue-router";
import Fim_event from "@/components/fim_event.vue";
import {chatUserTopApi} from "@/api/chat_api";
import {ElMessage} from "element-plus";
import {groupTopApi} from "@/api/group_api";

const route = useRoute()

interface sessionType {
  id: number
  title: string
  avatar: string
  date: string
  msgPreview: string
  isTop: boolean
  type: "user" | "group"
}

const data = reactive<listResponse<sessionType>>({
  list: [],
  count: 0,
})

async function getData() {
  const chatRes = await chatSessionListApi()
  const groupRes = await groupSessionListApi()


  const list: sessionType[] = []
  for (const li of chatRes.data.list) {
    list.push({
      id: li.userId,
      title: li.nickname,
      avatar: li.avatar,
      date: li.created_at,
      msgPreview: li.msgPreview,
      isTop: li.isTop,
      type: "user",
    })
  }
  for (const li of groupRes.data.list) {
    list.push({
      id: li.groupId,
      title: li.title,
      avatar: li.avatar,
      date: li.newMsgDate,
      msgPreview: li.newMsgPreview,
      isTop: li.isTop,
      type: "group",
    })
  }

  const newList = list.sort((a, b) => {
    // 如果a和b都有isTop呢
    if (a.isTop && b.isTop){
      return new Date(b.date).getTime() - new Date(a.date).getTime()
    }
    if (a.isTop) {
      // a如果有isTop，那么就要让它排到前面去
      return -1;
    }
    if (b.isTop) {
      // b如果有isTop，那就维持现状
      return 1;
    }
    return new Date(b.date).getTime() - new Date(a.date).getTime()
  })

  data.list = newList
}

function goRouter(record: sessionType) {
  if (record.type === "user") {
    router.push({
      name: "session_user_chat",
      params: {
        id: record.id,
      }
    })
    return
  }
  router.push({
    name: "session_group_chat",
    params: {
      id: record.id,
    }
  })
}

interface clsType {
  user?: boolean
  group?: boolean
  isTop?: boolean
  active?: boolean
}

const itemClass = computed(() => {
  return (item: sessionType) => {
    const cls: clsType = {}
    cls[item.type] = true
    cls.isTop = item.isTop
    const id = Number(route.params.id || '0')

    if (route.path.startsWith("/group") && item.type === "group") {
      if (id === item.id) {
        cls.active = true
      }
    }
    if (route.path.startsWith("/user") && item.type === "user") {
      if (id === item.id) {
        cls.active = true
      }
    }
    return cls
  }
})

getData()

function event() {
  // 先判断是不是用户聊天页面，如果是，就拿到用户id，然后在用户会话列表里面查
  // 如果不在，就去重新调一下列表
  if (route.name === "session_user_chat") {
    const userID = Number(route.params.id)
    const val = data.list.find((value) => value.type === "user" && value.id === userID)
    if (val === undefined) {
      // 没有找到
      getData()
    }
  }
}

async function userChatTop(isTop: boolean, friendID: number) {
  let res = await chatUserTopApi(friendID)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  if (isTop) {
    ElMessage.success("置顶用户成功")
  } else {
    ElMessage.success("取消置顶用户成功")
  }
  getData()

}

async function groupChatTop(isTop: boolean, groupID: number) {
  let res = await groupTopApi(groupID, isTop)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  if (isTop) {
    ElMessage.success("置顶群成功")
  } else {
    ElMessage.success("取消置顶群成功")
  }
  getData()

}


</script>

<template>
  <div class="session_view">
    <fim_event event-key="sessionList" @event="event"></fim_event>
    <div class="session_slide">
      <fim_slide_head></fim_slide_head>
      <el-scrollbar style="height: calc(100% - 40px)">
        <div class="session_list">
          <template v-for="item in data.list">
            <el-dropdown trigger="contextmenu">
              <div class="item" @click="goRouter(item)" :class="itemClass(item)">
                <fim_avatar shape="circle" :src="item.avatar"></fim_avatar>
                <div class="info">
                  <el-text class="title" truncated style="max-width: 6rem">{{ item.title }}</el-text>
                  <el-text class="msg" truncated>{{ item.msgPreview }}</el-text>
                </div>
                <div class="date">{{ relativeCurrentTime(item.date) }}</div>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item v-if="item.type === 'user' && !item.isTop" @click="userChatTop(true, item.id)">
                    置顶用户
                  </el-dropdown-item>
                  <el-dropdown-item v-if="item.type === 'group' && !item.isTop" @click="groupChatTop(true, item.id)">
                    置顶群
                  </el-dropdown-item>
                  <el-dropdown-item v-if="item.type === 'user' && item.isTop" @click="userChatTop(false, item.id)">
                    取消置顶
                  </el-dropdown-item>
                  <el-dropdown-item v-if="item.type === 'group' && item.isTop" @click="groupChatTop(false, item.id)">
                    取消置顶
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>

        </div>
      </el-scrollbar>

    </div>
    <div class="session_main">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss">
.session_view {
  width: 100%;
  display: flex;
  height: 100%;

  .session_slide {
    width: 200px;
    height: 100%;
    border-right: 1px solid var(--border_color);

    .session_list {
      width: 100%;

      .el-dropdown {
        width: 100%;
      }

      .item {
        width: 100%;
        display: flex;
        padding: 10px;
        cursor: pointer;
        justify-content: space-between;
        position: relative;

        &:hover {
          background-color: var(--item_hover);
        }

        &.active {
          background-color: var(--item_hover);
        }

        &.group {
          .title {
            color: red;
          }
        }

        &.isTop {
          &::before {
            content: "";
            display: block;
            position: absolute;
            left: 5px;
            top: 5px;
            width: 0;
            height: 0;
            border: 5px solid transparent;
            border-left-color: var(--text_color);
            border-top-color: var(--text_color);
          }
        }

        .fim_avatar {
          width: 40px;
        }

        .info {
          width: calc(100% - 45px);
          display: flex;
          flex-direction: column;

          .el-text {
            align-self: start;
          }

          .title {
            height: 1rem;
          }

          .msg {
            margin-top: 5px;
          }
        }

        .date {
          position: absolute;
          right: 10px;
          top: 10px;
          font-size: 12px;
          color: var(--text_color);
        }
      }
    }


  }

  .session_main {
    width: calc(100% - 200px);
    height: 100%;
  }
}
</style>