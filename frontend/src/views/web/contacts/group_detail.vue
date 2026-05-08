<script setup lang="ts">
import {useRoute} from "vue-router";
import {reactive} from "vue";
import type {groupInfoType} from "@/api/group_api";
import Fim_avatar from "@/components/fim_avatar.vue";
import {ElMessage} from "element-plus";
import {watch} from "vue";
import {groupDetailApi, groupMemberExitApi} from "@/api/group_api";
import router from "@/router";
import {useStore} from "@/stores";

const store = useStore()
const route = useRoute()

const data = reactive<groupInfoType>({
  "groupId": 0,
  "title": "",
  "abstract": "",
  "memberCount": 0,
  "memberOnlineCount": 0,
  "avatar": "",
  "creator": {
    "userId": 0,
    "avatart": "",
    "nickname": ""
  },
  "adminList": [],
  "role": 0,
  "isProhibition": false,
  prohibitionTime: undefined,
  isSearch: false,
  isInvite: false,
  isTemporarySession: false,
})


async function getData() {
  let res = await groupDetailApi(Number(route.params.id))
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  Object.assign(data, res.data)
}

watch(() => route.params.id, () => {
  getData()
}, {immediate: true})

function goGroupChat() {
  router.push({
    name: "session_group_chat",
    params: {
      id: data.groupId,
    }
  })
}

function goGroupSettings() {
  router.push({
    name: "group_settings",
    params: {
      id: data.groupId,
    }
  })
}

async function memberExit() {
  let res = await groupMemberExitApi({
    memberId: store.userInfo.userID,
    id: data.groupId,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("退出群聊成功")
  router.push({
    name: "contacts",
  })
}

</script>

<template>
  <div class="group_detail_view">
    <div class="top">
      <div class="left">
        <div class="nickname">
          <el-text style="max-width: 30rem" truncated>{{ data.title }}</el-text>
        </div>
        <div class="item">
          <span class="label">ID</span>
          <span class="val">{{ data.groupId }}</span>
        </div>
        <div class="item">
          <span class="label">人数</span>
          <span class="val">{{ data.memberOnlineCount }}/{{ data.memberCount }}</span>
        </div>
        <div class="item">
          <span class="label">简介</span>
          <span class="val">{{ data.abstract }}</span>
        </div>
      </div>
      <div class="right">
        <div class="more">
          <el-dropdown trigger="click">
            <i class="iconfont icon-more"></i>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goGroupChat">进入群聊</el-dropdown-item>
                <el-dropdown-item v-if="data.role !== 3" @click="goGroupSettings">群设置</el-dropdown-item>
                <el-dropdown-item @click="memberExit" style="color: red" v-if="data.role !== 1">退出群聊</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="avatar">
          <fim_avatar :src="data.avatar" shape="square" :size="60"></fim_avatar>
        </div>
      </div>
    </div>
    <div class="bottom">
      <el-button type="primary" @click="goGroupChat">进入群聊</el-button>
    </div>
  </div>
</template>

<style lang="scss">
.group_detail_view {
  padding: 40px 20px;

  .top {
    border-bottom: 1px solid var(--border_color);
    display: flex;
    padding: 0 0 20px 0;
    margin: 0 0 20px 0;


    .left {
      width: 70%;

      .nickname {
        font-size: 22px;
        display: flex;
        align-items: center;
        color: var(--text_color);

        > span {
          font-size: 22px;
        }
      }

      .item {
        color: var(--text_color);
        margin-top: 10px;

        .label {
          font-weight: 600;

          &::after {
            content: ":";
            display: inline-block;
            margin-right: 5px;
          }
        }
      }
    }

    .right {
      width: 30%;
      text-align: right;

      .more {
        i {
          font-size: 22px;
          cursor: pointer;
        }
      }
    }
  }

  .bottom {
    text-align: center;
  }
}
</style>