<script setup lang="ts">
import {friendRemoveApi, friendUpdateNicknameApi} from "@/api/user_api";
import {nextTick, reactive, ref, watch} from "vue";
import type {friendType} from "@/api/user_api";
import {useRoute} from "vue-router";
import {friendDetailApi} from "@/api/user_api";
import {ElMessage} from "element-plus";
import router from "@/router";
import {eventStore} from "@/stores/event_store";
import Fim_avatar from "@/components/fim_avatar.vue";
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



</script>

<template>
  <div class="chat_user_detail_view">
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
    <el-form>
      <el-form-item label="头像">
        <fim_avatar :src="userDetail.avatar"></fim_avatar>
      </el-form-item>
      <el-form-item label="用户号">
       {{ userDetail.userID }}
      </el-form-item>
      <el-form-item label="昵称">
        {{ userDetail.nickname }}
      </el-form-item>
      <el-form-item label="备注">
        <fim_friend_notice :friend-id="userDetail.userID" v-model="userDetail.notice"></fim_friend_notice>
      </el-form-item>
      <el-form-item label="简介">
        {{ userDetail.abstract }}
      </el-form-item>
    </el-form>
  </div>
</template>

<style lang="scss">
.chat_user_detail_view {
  padding: 10px;
  .el-form-item{
    margin-bottom: 8px;
  }

}
</style>