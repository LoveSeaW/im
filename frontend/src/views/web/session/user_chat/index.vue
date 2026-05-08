<script setup lang="ts">
import router from "@/router";
import {useRoute} from "vue-router";
import {useStore} from "@/stores";
import {ElMessage} from "element-plus";
import {reactive, watch} from "vue";
import type {friendType} from "@/api/user_api";
import {friendDetailApi} from "@/api/user_api";

const route = useRoute()
const store = useStore()

const data = reactive<friendType>({
  userID: 0,
  nickname: "",
  abstract: "",
  avatar: "",
  notice: "",
  isOnline: false,
})

async function getData() {
  let res = await friendDetailApi({friendID: Number(route.params.id)})
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  Object.assign(data, res.data)
}

watch(() => route.params, () => {
  getData()
}, {immediate: true, deep: true})


function goGroupRoute(name: string) {
  router.push({
    name: name,
  })
}

</script>

<template>
  <div class="user_chat_view">
    <div class="user_chat_head">
      <div class="title">
        {{ data.nickname }}
      </div>
      <div class="icon">
        <el-dropdown trigger="click">
          <i class="iconfont icon-more"></i>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="goGroupRoute('session_user_info')">用户资料</el-dropdown-item>
              <el-dropdown-item style="color:red;">
                删除用户
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    <div class="user_chat_main">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss">
.user_chat_view {
  height: 100%;

  .user_chat_head {
    height: 40px;
    border-bottom: 1px solid var(--border_color);
    display: flex;
    justify-content: space-between;
    padding: 0 10px;
    align-items: center;
    color: var(--text_color);

    .icon {
      cursor: pointer;
    }
  }

  .user_chat_main {
    height: 100%;
  }
}
</style>