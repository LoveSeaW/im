<script setup lang="ts">
import {ArrowDown} from "@element-plus/icons-vue";
import {reactive} from "vue";
import type {listResponse} from "@/api";
import type {userValidType} from "@/api/user_api";
import {userValidListApi} from "@/api/user_api";
import {userValidStatusApi} from "@/api/user_api";
import {ElMessage} from "element-plus";
import {CloseBold} from "@element-plus/icons-vue";
import {useStore} from "@/stores";
import {groupValidListApi, type groupValidType} from "@/api/group_api";
import {groupValidStatusApi} from "@/api/group_api";
import Valid_detail, {type ValidType} from "@/components/valid_detail.vue";

const store = useStore()
const data = reactive<listResponse<userValidType>>({
  list: [],
  count: 0
})

async function getData() {
  let res = await userValidListApi()
  data.list = res.data.list || []
  data.count = res.data.count
}

async function userValidStatus(id: number, status: number) {
  let res = await userValidStatusApi({
    status,
    verifyId: id,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("成功")
  getData()
}

const activeUserValidData = reactive<userValidType>({
  userID: 0,
  nickname: "",
  avatar: "",
  additionalMessages: "",
  verificationQuestion: undefined,
  status: 0,
  sendStatus: 0,
  revStatus: 0,
  verification: 0,
  id: 0,
  flag: "send",
  createdAt: "",
})

function checkValid(record: userValidType) {
  Object.assign(activeUserValidData, record)

  if (record.flag === "send") {
    activeData.sendName = store.userInfo.nickname
    activeData.sendAvatar = store.userInfo.avatar
    activeData.revName = record.nickname
    activeData.revAvatar = record.avatar
  } else {
    activeData.sendName = record.nickname
    activeData.sendAvatar = record.avatar
    activeData.revName = store.userInfo.nickname
    activeData.revAvatar = store.userInfo.avatar
  }
  activeData.createdAt = record.createdAt
  activeData.additionalMessages = record.additionalMessages
  activeData.verificationQuestion = record.verificationQuestion
}


getData()


const groupData = reactive<listResponse<groupValidType>>({
  list: [],
  count: 0,
})

async function getGroupData() {
  let res = await groupValidListApi()
  groupData.list = res.data.list || []
  groupData.count = res.data.count
}

function checkGroupValid(record: groupValidType) {
  Object.assign(activeGroupValidData, record)
  activeData.sendName = record.userNickname
  activeData.sendAvatar = record.userAvatar
  activeData.revName = record.title
  activeData.revAvatar = record.title[0]
  activeData.createdAt = record.createdAt
  activeData.additionalMessages = record.additionalMessages
  activeData.verificationQuestion = record.verificationQuestion
}

async function groupValidStatus(id: number, status: number) {
  let res = await groupValidStatusApi({
    status,
    validId: id,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("成功")
  getGroupData()
}

const activeData = reactive<ValidType>({
  sendAvatar: "",
  sendName: "",
  revAvatar: "",
  revName: "",
  createdAt: "",
  additionalMessages: "",
  verificationQuestion: undefined,
})
const activeGroupValidData = reactive<groupValidType>({
  id: 0,
  grouId: 0,
  userId: 0,
  userNickname: "",
  userAvatar: "",
  status: 0,
  additionalMessages: "",
  verificationQuestion: undefined,
  title: "",
  createdAt: "",
  type: 0,
})


getGroupData()

</script>

<template>
  <div class="notice_view">
    <div class="notice_slide">
      <el-scrollbar height="100%">
        <el-menu :default-openeds="['1', '2']">
          <el-sub-menu index="1">
            <template #title>
              <span>好友验证</span>
            </template>
            <div class="valid_list">
              <div class="item" :class="{active: activeUserValidData.id === item.id}" @click="checkValid(item)"
                   v-for="item in data.list">
                <el-avatar class="left" :src="item.avatar" :size="40"></el-avatar>
                <div class="center">
                  <el-text class="nickname" truncated style="max-width: 10rem">{{ item.nickname }}</el-text>
                  <el-text class="msg" truncated style="max-width: 10rem">附加消息：{{
                      item.additionalMessages
                    }}
                  </el-text>
                </div>
                <div class="right">
                  <template v-if="item.flag === 'rev'">
                    <el-button-group v-if="item.revStatus === 0">
                      <el-button type="primary" @click="userValidStatus(item.id, 1)" size="small">同意</el-button>
                      <el-dropdown>
                        <el-button style="width: 20px;" :icon="ArrowDown as unknown" size="small"></el-button>
                        <template #dropdown>
                          <el-dropdown-item @click="userValidStatus(item.id, 2)">拒绝</el-dropdown-item>
                        </template>
                      </el-dropdown>
                    </el-button-group>
                    <el-button v-if="item.revStatus === 0" style="margin-left: 5px" @click="userValidStatus(item.id, 3)"
                               size="small">忽略
                    </el-button>
                    <span class="status_1" v-if="item.revStatus === 1">已同意</span>
                    <span class="status_2" v-if="item.revStatus === 2">已拒绝</span>
                    <span class="status_3" v-if="item.revStatus === 3">已忽略</span>
                  </template>

                  <template v-if="item.flag === 'send'">
                    <span v-if="item.revStatus === 0">等待对方验证</span>
                    <span class="status_1" v-if="item.revStatus === 1">对方已同意</span>
                    <span class="status_2" v-if="item.revStatus === 2">对方已拒绝</span>
                    <span class="status_3" v-if="item.revStatus === 3">对方已忽略</span>
                  </template>

                  <div class="icon" v-if="item.revStatus !== 0">
                    <el-icon size="14" @click="userValidStatus(item.id, 4)">
                      <CloseBold/>
                    </el-icon>
                  </div>
                </div>

              </div>
            </div>

          </el-sub-menu>
          <el-sub-menu index="2">
            <template #title>
              <span>群系统消息</span>
            </template>
            <div class="valid_list group_valid_list">
              <div class="item" :class="{active: activeGroupValidData.id === item.id}" @click="checkGroupValid(item)"
                   v-for="item in groupData.list">
                <el-avatar class="left" :src="item.userAvatar" :size="40"></el-avatar>
                <div class="center">
                  <el-text class="nickname" truncated style="max-width: 10rem">
                    <template v-if="item.type === 1">{{ item.userNickname }}申请加入群 {{ item.title }}</template>
                    <template v-if="item.type === 2">{{ item.userNickname }}退出群 {{ item.title }}</template>
                  </el-text>
                  <el-text class="msg" truncated style="max-width: 10rem" v-if="item.type === 1">附加消息：{{
                      item.additionalMessages
                    }}
                  </el-text>
                </div>
                <div class="right">
                  <template v-if="item.type === 1">
                    <template v-if="item.status === 0">
                      <el-button-group>
                        <el-button type="primary" @click="groupValidStatus(item.id, 1)" size="small">同意</el-button>
                        <el-dropdown>
                          <el-button style="width: 20px;" :icon="ArrowDown as unknown" size="small"></el-button>
                          <template #dropdown>
                            <el-dropdown-item @click="groupValidStatus(item.id, 2)">拒绝</el-dropdown-item>
                          </template>
                        </el-dropdown>
                      </el-button-group>
                      <el-button style="margin-left: 5px" @click="groupValidStatus(item.id, 3)" size="small">忽略
                      </el-button>
                    </template>
                    <template v-if="item.status === 1">
                      <span class="status_1">已同意</span>
                    </template>
                    <template v-if="item.status === 2">
                      <span class="status_2">已拒绝</span>
                    </template>
                    <template v-if="item.status === 3">
                      <span class="status_3">已忽略</span>
                    </template>
                  </template>

                  <div class="icon" v-if="item.status !== 0">
                    <el-icon size="14" @click="groupValidStatus(item.id, 4)">
                      <CloseBold/>
                    </el-icon>
                  </div>
                </div>

              </div>
            </div>
          </el-sub-menu>
        </el-menu>
      </el-scrollbar>
    </div>
    <div class="notice_main">
      <valid_detail v-if="activeUserValidData.id || activeGroupValidData.id" :data="activeData"></valid_detail>
    </div>
  </div>
</template>

<style lang="scss">
.notice_view {
  width: 100%;
  display: flex;
  height: 100%;

  .notice_slide {
    width: 340px;
    border-right: 1px solid var(--border_color);
    height: 100%;


    .el-menu {
      border-right: none;

      .el-sub-menu__title {
        height: 40px;
        font-weight: 600;
        padding: 0 10px;
      }


      .valid_list {
        width: 100%;

        .item {
          height: 60px;
          display: flex;
          padding: 10px;
          align-items: center;
          cursor: pointer;
          position: relative;

          &:hover {
            background-color: var(--item_hover);

            .icon {
              opacity: 1 !important;
            }
          }

          &.active {
            background-color: var(--item_hover);
          }

          .left {
            flex-shrink: 0;
          }

          .right {
            flex-shrink: 0;
            font-size: 14px;
            color: var(--text_color);

            .icon {
              position: absolute;
              right: 0;
              top: 0;
              opacity: 0;
              transition: all .3s;
            }

            .status_1 {
              color: #529b2e;
            }

            .status_2 {
              color: #cb1b2f;
            }

            .status_3 {
              color: #6e6e6e;
            }
          }

          .center {
            flex-shrink: 1;
            width: 100%;
            padding: 0 5px;
            display: flex;
            flex-direction: column;

            .el-text {
              align-self: start;
            }
          }

        }
      }
    }
  }

  .notice_main {
    width: calc(100% - 340px);


  }
}
</style>