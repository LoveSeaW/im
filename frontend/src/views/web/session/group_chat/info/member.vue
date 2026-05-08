<script setup lang="ts">
import {useStore} from "@/stores";
import {groupMemberListApi, type groupMemberListRequest, type groupMemberType} from "@/api/group_api";
import {reactive} from "vue";
import type {listResponse} from "@/api";
import {dateTimeFormat, relativeCurrentTime} from "@/utils/date";
import Fim_avatar from "@/components/fim_avatar.vue";
import {groupMemberRoleUpdateApi} from "@/api/group_api";
import {ElMessage} from "element-plus";
import {groupMemberExitApi} from "@/api/group_api";
import {Edit} from "@element-plus/icons-vue";
import {ElMessageBox} from "element-plus";
import {groupMemberNicknameApi} from "@/api/group_api";

const store = useStore()


const params = reactive<groupMemberListRequest>({
  id: store.groupData.groupId,
  sort: "role asc"
})
const data = reactive<listResponse<groupMemberType>>({
  list: [],
  count: 0,
})

async function getData() {
  let res = await groupMemberListApi(params)
  data.list = res.data.list || []
  data.count = res.data.count
  // console.log(data.list)
}

async function updateMemberRole(record: groupMemberType, role: number) {
  let res = await groupMemberRoleUpdateApi({
    id: store.groupData.groupId,
    memberId: record.userId,
    role: role
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("成员角色更新成功")
  getData()
}


async function groupMemberExit(record: groupMemberType) {
  let res = await groupMemberExitApi({
    id: store.groupData.groupId,
    memberId: record.userId,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("成员移出成功")
  getData()
}

function showMemberNickname(record: groupMemberType) {
  ElMessageBox.prompt('群成员备注', '修改群成员备注', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputValue: record.memberNickname,
  })
      .then(async ({value}) => {
        let res = await groupMemberNicknameApi({
          id: store.groupData.groupId,
          memberId: record.userId,
          nickname: value,
        })
        if (res.code) {
          ElMessage.error(res.msg)
          return
        }
        ElMessage.success("修改群成员备注成功")
        getData()
      })
      .catch(() => {
      })
}

getData()

</script>

<template>
  <div class="member_view">
    <el-table :data="data.list" style="width: 100%">
      <el-table-column label="成员" width="150">
        <template #default="{row}:{row: groupMemberType}">
          <div class="info">
            <div class="icon">
              <i class="iconfont icon-qunzhu" v-if="row.role === 1"></i>
              <i class="iconfont icon-qunguanliyuan" v-if="row.role === 2"></i>
            </div>
            <fim_avatar :size="30" :src="row.avatar"></fim_avatar>
            <el-text :title="row.userNickname" truncated style="max-width: 5rem" class="nickname">{{ row.userNickname }}</el-text>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="群昵称" width="100">
        <template #default="{row}:{row: groupMemberType}">
          <div class="col_group_nickname">
            <el-text :title="row.memberNickname" truncated style="max-width: 6rem">{{ row.memberNickname }}</el-text>
            <el-icon @click="showMemberNickname(row)" v-if="
            store.groupData.role === 1 ||
             (store.groupData.role === 2 && (row.role === 3 || (row.role === 2 && store.userInfo.userID === row.userId)))||
              (store.groupData.role === 3 && store.userInfo.userID === row.userId)" size="18">
              <Edit></Edit>
            </el-icon>
          </div>

        </template>
      </el-table-column>
      <el-table-column label="加群时间" width="120">
        <template #default="scope">
          <span>{{ relativeCurrentTime(scope.row.createdAt) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="最后发言" width="120">
        <template #default="{row}:{row: groupMemberType}">
          <span v-if="row.newMsgDate">{{ relativeCurrentTime(row.newMsgDate) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="210">
        <template #default="{row}:{row: groupMemberType}">
          <div class="btn_list" style="display: flex; flex-shrink: 0">
            <el-button @click="groupMemberExit(row)" type="danger" size="small"
                       v-if="(store.groupData.role === 1 && row.role !== 1) || (store.groupData.role === 2 && row.role === 3)">
              踢出群聊
            </el-button>
            <el-button type="primary" size="small" v-if="(store.groupData.role === 1 && row.role === 3)"
                       @click="updateMemberRole(row, 2)">升为管理员
            </el-button>
            <el-button type="primary" size="small" v-if="(store.groupData.role === 1 && row.role === 2)"
                       @click="updateMemberRole(row, 3)">降为群成员
            </el-button>
          </div>

        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style lang="scss">
.member_view {
  .info {
    display: flex;
    align-items: center;

    .icon {
      margin-right: 5px;
      font-size: 12px;
    }

    .nickname {
      margin-left: 5px;
    }
  }

  .col_group_nickname {
    display: flex;
    align-items: center;

    .el-icon {
      cursor: pointer;
    }
  }
}
</style>