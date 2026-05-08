<script setup lang="ts">
import router from "@/router";
import {useRoute} from "vue-router";
import {useStore} from "@/stores";
import {groupMemberExitApi} from "@/api/group_api";
import {ElMessage} from "element-plus";
import {watch} from "vue";

const route = useRoute()
const store = useStore()


watch(() => route.params, () => {
  store.getGroupData(Number(route.params.id))
}, {immediate: true, deep: true})

function goGroupRoute(name: string) {
  router.push({
    name: name,
  })
}


async function memberExit() {
  let res = await groupMemberExitApi({
    memberId: store.userInfo.userID,
    id: store.groupData.groupId,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("退出群聊成功")
  router.push({
    name: "session",
  })
}

</script>

<template>
  <div class="group_chat_view">
    <div class="group_chat_head">
      <div class="title">
        {{ store.groupData.title }}
      </div>
      <div class="icon">
        <el-dropdown trigger="click">
          <i class="iconfont icon-more"></i>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="goGroupRoute('session_group_chat')">群对话</el-dropdown-item>
              <el-dropdown-item @click="goGroupRoute('group_information')">群资料</el-dropdown-item>
              <el-dropdown-item @click="goGroupRoute('group_member')" v-if="store.groupData.role !== 3">群成员
              </el-dropdown-item>
              <el-dropdown-item @click="goGroupRoute('group_settings')" v-if="store.groupData.role !== 3">群设置
              </el-dropdown-item>
              <el-dropdown-item style="color:red;" @click="memberExit" divided v-if="store.groupData.role !== 1">
                退出群聊
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
    <div class="group_chat_main">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss">
.group_chat_view {
  height: 100%;

  .group_chat_head {
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

  .group_chat_main {
    height: calc(100% - 40px);
  }
}
</style>