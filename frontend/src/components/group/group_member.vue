<script setup lang="ts">
import {computed, reactive, ref, watch} from "vue";
import type {listResponse} from "@/api";
import {
  type groupFriendResponse, groupMemberAddApi,
  groupMemberListApi,
  type groupMemberListRequest, groupMemberProhibitionApi,
  type groupMemberRequest,
  type groupMemberType, groupMyFriendListApi
} from "@/api/group_api";
import {ElMessage} from "element-plus";
import {showAddUserModal} from "@/components/add_user_modal";
import {useStore} from "@/stores";
import {useRouter} from "vue-router";
import Fim_avatar from "@/components/fim_avatar.vue";

const router = useRouter()
const store = useStore()

interface Props {
  groupId: number
}

const props = defineProps<Props>()


const data = reactive<listResponse<groupMemberType>>({
  list: [],
  count: 0,
})

const groupMemberListParams = reactive<groupMemberListRequest>({
  id: props.groupId, // 注意，这里去切换群的适合，这个值还是之前的值
  sort: "role asc",
  key: ""
})

async function getGroupMemberList() {
  groupMemberListParams.id = props.groupId
  let res = await groupMemberListApi(groupMemberListParams)
  data.list = res.data.list || []
  data.count = res.data.count
}


watch(() => props.groupId, () => {
  console.log(props.groupId)
  getGroupMemberList()
}, {immediate: true})


const visible = ref(false)

const groupMemberReq = reactive<groupMemberRequest>({
  id: props.groupId,
  memberIdList: []
})
const myFriendList = ref<groupFriendResponse[]>([])

async function showAddMemberGroup() {
  let res = await groupMyFriendListApi(props.groupId)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  myFriendList.value = res.data.list
  visible.value = true
}

const useFriendList = computed(() => {
  const list: groupFriendResponse[] = []
  for (const argument of myFriendList.value) {
    if (groupMemberReq.memberIdList.includes(argument.userId)) {
      list.push(argument)
    }
  }
  return list
})

function unUseFriend(userID: number) {
  const index = groupMemberReq.memberIdList.findIndex((value) => value === userID)
  if (index === -1) {
    return
  }
  groupMemberReq.memberIdList.splice(index, 1)
}

async function groupMemberAdd() {
  let res = await groupMemberAddApi(groupMemberReq)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("邀请成功")
  getGroupMemberList()
  groupMemberClear()
}

async function groupMemberClear() {
  groupMemberReq.memberIdList = []
  visible.value = false
}

const groupOnlineMemberCount = computed(() => {
  return data.list.filter((value) => value.isOnline).length
})

const searchIpt = ref()

function searchMember() {
  getGroupMemberList()
}

function searchIptFocus() {
  searchIpt.value.focus()
}

function addFriend(item: groupMemberType) {
  showAddUserModal({
    userID: item.userId,
    nickname: item.userNickname,
    abstract: "",
    avatar: item.avatar,
    isFriend: false,
  })
}

async function groupMemberProhibition(item: groupMemberType, prohibitionTime?: number) {
  let res = await groupMemberProhibitionApi({
    groupId: props.groupId,
    memberId: item.userId,
    prohibitionTime: prohibitionTime,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("禁言设置成功")
}

function goRoute(item: groupMemberType, name: string) {
  router.push({
    name: name,
    params: {
      id: item.userId,
    }
  })
}


</script>

<template>
  <div class="fim_group_member">
    <el-dialog width="30%" title="邀请进群" v-model="visible">
      <div class="add_member_group_dialog">
        <div class="dialog_left">
          <div class="title">我的好友</div>
          <el-scrollbar style="max-height: 350px">
            <el-checkbox-group v-model="groupMemberReq.memberIdList" class="member_list">
              <el-checkbox class="item" v-for="item in myFriendList" :value="item.userId" :disabled="item.isInGroup">
                <fim_avatar :size="20" :src="item.avatar"></fim_avatar>
                <el-text truncated>{{ item.nickname }}</el-text>
              </el-checkbox>
            </el-checkbox-group>
          </el-scrollbar>
        </div>
        <div class="dialog_right">
          <div class="title">
            已选联系人
          </div>
          <div class="member_list">
            <el-scrollbar style="max-height: 350px">
              <div class="item item_check" v-for="item in useFriendList">
                <fim_avatar :size="20" :src="item.avatar"></fim_avatar>
                <el-text truncated>{{ item.nickname }}</el-text>
                <i @click="unUseFriend(item.userId)" class="iconfont icon-quxiao"></i>
              </div>
            </el-scrollbar>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button type="primary" @click="groupMemberAdd">确定</el-button>
        <el-button @click="groupMemberClear">取消</el-button>
      </template>
    </el-dialog>
    <div class="head">
      <span>群成员<span>{{ groupOnlineMemberCount }}/{{ data.list.length }}</span></span>
      <span>
            <i @click="showAddMemberGroup" class="iconfont icon-add-group"></i>
            <i @click="searchIptFocus" class="iconfont icon-sousuo"></i>
          </span>
    </div>
    <div class="search">
      <el-input ref="searchIpt" v-model="groupMemberListParams.key" @keydown.enter="searchMember" size="small"
                placeholder="搜索群成员"></el-input>
    </div>
    <div class="member_list">
      <el-scrollbar style="height: 100%">
        <el-dropdown trigger="contextmenu" v-for="item in data.list">
          <div class="item">
            <fim_avatar :online="item.isOnline" :src="item.avatar"></fim_avatar>
            <el-text class="name" truncated style="max-width: 5rem">{{ item.userNickname }}</el-text>
            <div class="icon">
              <i class="iconfont icon-qunzhu" v-if="item.role === 1"></i>
              <i class="iconfont icon-qunguanliyuan" v-else-if="item.role === 2"></i>
              <i v-if="item.prohibitionTime" class="iconfont icon-jinyan"></i>
            </div>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <template v-if="item.isFriend">
                <el-dropdown-item @click="goRoute(item, 'session_user_chat')">私聊用户</el-dropdown-item>
                <el-dropdown-item @click="goRoute(item, 'user_detail')">用户资料</el-dropdown-item>
              </template>
              <template v-else>
                <el-dropdown-item @click="addFriend(item)">添加好友</el-dropdown-item>
              </template>

              <el-dropdown-item
                  v-if="!item.prohibitionTime && ((store.groupData.role === 1 && item.role !== 1) || (store.groupData.role === 2 && item.role === 3))">
                <el-dropdown trigger="hover">
                  用户禁言
                  <template #dropdown>
                    <el-dropdown-item @click="groupMemberProhibition(item, 5)">五分钟</el-dropdown-item>
                    <el-dropdown-item @click="groupMemberProhibition(item, 30)">30分钟</el-dropdown-item>
                    <el-dropdown-item @click="groupMemberProhibition(item, 60)">1小时</el-dropdown-item>
                  </template>
                </el-dropdown>
              </el-dropdown-item>
              <el-dropdown-item @click="groupMemberProhibition(item)"
                                v-if="item.prohibitionTime &&((store.groupData.role === 1 && item.role !== 1) || (store.groupData.role === 2 && item.role === 3))">
                取消禁言
              </el-dropdown-item>
              <el-dropdown-item style="color: red">踢出群聊</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

      </el-scrollbar>
    </div>
  </div>
</template>

<style lang="scss">
.fim_group_member {

  height: calc(100% - 150px);

  .head {
    padding: 10px;
    font-size: 14px;
    color: var(--text_color);
    display: flex;
    justify-content: space-between;

    i {
      cursor: pointer;

      &:hover {
        color: var(--active);
      }
    }
  }

  .search {
    padding: 0 10px;
  }

  .member_list {
    height: calc(100% - 49px);

    .el-dropdown {
      width: 100%;

    }

    .item {
      padding: 5px 10px;
      display: flex;
      align-items: center;
      position: relative;
      width: 100%;

      .name {
        margin-left: 5px;
      }

      .icon {
        position: absolute;
        right: 15px;
      }

      &:hover {
        background-color: var(--item_hover);
      }

      .el-text {
        height: 16px;
      }
    }
  }
}
</style>