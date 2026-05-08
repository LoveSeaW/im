<script setup lang="ts">
import {friendRemoveApi, friendUpdateNicknameApi} from "@/api/user_api";
import {nextTick, reactive, ref, watch} from "vue";
import type {friendType} from "@/api/user_api";
import {useRoute} from "vue-router";
import {friendDetailApi} from "@/api/user_api";
import {ElMessage} from "element-plus";
import router from "@/router";
import {eventStore} from "@/stores/event_store";
import Fim_friend_notice from "@/components/fim_friend_notice.vue";

const es = eventStore()
const route = useRoute()

const userDetail = reactive<friendType>({
  nickname: "",
  notice: "",
  userID: 0,
  abstract: "",
  avatar: "",
  isOnline: false,
})

async function getData() {
  let res = await friendDetailApi({friendID: Number(route.params.id)})
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  Object.assign(userDetail, res.data)
}

watch(() => route.params.id, () => {
  getData()
}, {immediate: true})


const visible = ref(false)
const updateNoticeRef = ref()

function updateNoticeShow() {
  visible.value = true
  setTimeout(() => {
    updateNoticeRef.value.focus()
  }, 100)

}

async function updateFriendNickname() {
  let res = await friendUpdateNicknameApi(userDetail.userID, userDetail.notice)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("修改成功")
  visible.value = false
  return
}

function goUserChat() {
  router.push({
    name: "session_user_chat",
    params: {
      id: userDetail.userID,
    }
  })
}

function goUserInfo() {
  router.push({
    name: "session_user_info",
    params: {
      id: userDetail.userID,
    }
  })
}

async function friendRemove() {
  let res = await friendRemoveApi(userDetail.userID)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("删除好友成功")
  es.setEventKey("friendList")
  router.push({
    name: "welcome"
  })

}

</script>

<template>
  <div class="user_detail_view">
    <el-dialog v-model="visible" width="30%" title="修改好友备注">
      <el-input ref="updateNoticeRef" placeholder="备注" v-model="userDetail.notice"></el-input>
      <template #footer>
        <div class="dialog-footer">
          <el-button type="primary" @click="updateFriendNickname">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
    <div class="top">
      <div class="left">
        <div class="nickname">
          <el-text style="max-width: 10rem" truncated>{{ userDetail.nickname }}</el-text>
          （
          <el-text style="max-width: 15rem" truncated>
            {{ userDetail.notice === "" ? "-" : userDetail.notice }}
          </el-text>
          ）
        </div>
        <div class="item">
          <span class="label">ID</span>
          <span class="val">{{ userDetail.userID }}</span>
        </div>
        <div class="item">
          <span class="label">备注</span>
          <div class="val">
            <fim_friend_notice :friend-id="userDetail.userID" v-model="userDetail.notice"></fim_friend_notice>
          </div>
        </div>
        <div class="item">
          <span class="label">简介</span>
          <span class="val">{{ userDetail.abstract }}</span>
        </div>
      </div>
      <div class="right">
        <div class="more">
          <el-dropdown trigger="click">
            <i class="iconfont icon-more"></i>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goUserChat">发送消息</el-dropdown-item>
                <el-dropdown-item @click="goUserInfo">用户资料</el-dropdown-item>
                <el-dropdown-item @click="updateNoticeShow">修改备注</el-dropdown-item>
                <el-dropdown-item @click="friendRemove" style="color: red">删除用户</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="avatar">
          <el-avatar :size="60" shape="square" :src="userDetail.avatar"></el-avatar>
        </div>
      </div>
    </div>
    <div class="bottom">
      <el-button type="primary" @click="goUserChat">发送消息</el-button>
    </div>
  </div>
</template>

<style lang="scss">
.user_detail_view {
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
        display: flex;
        align-items: center;

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