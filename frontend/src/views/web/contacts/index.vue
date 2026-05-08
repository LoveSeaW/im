<script setup lang="ts">

import {reactive, watch} from "vue";
import type {listResponse} from "@/api";
import type {friendType} from "@/api/user_api";
import {friendListApi} from "@/api/user_api";
import {computed} from "vue";
import type {groupType} from "@/api/group_api";
import {groupMyListApi} from "@/api/group_api";
import Group_list from "@/components/group_list.vue";
import router from "@/router";
import {useRoute} from "vue-router";
import Fim_slide_head from "@/components/fim_slide_head.vue";
import {useStore} from "@/stores";
import Fim_event from "@/components/fim_event.vue";

const store = useStore()
const route = useRoute()
const friendData = reactive<listResponse<friendType>>({
  list: [],
  count: 0
})

async function getFriend() {
  let res = await friendListApi({limit: -1})
  friendData.list = res.data.list || []
  friendData.count = res.data.count
}

const friendOnlineCount = computed(() => {
  return friendData.list.filter((item) => item.isOnline).length
})

getFriend()

function checkFriend(record: friendType) {
  router.push({
    name: "user_detail",
    params: {
      id: record.userID,
    }
  })
}

const myJoinGroupData = reactive<listResponse<groupType>>({
  list: [],
  count: 0
})

async function getGroup() {
  let res = await groupMyListApi(2)
  myJoinGroupData.list = res.data.list || []
  myJoinGroupData.count = res.data.count
}


const myCreateGroupData = reactive<listResponse<groupType>>({
  list: [],
  count: 0
})

async function getMyCreateGroup() {
  let res = await groupMyListApi(1)
  myCreateGroupData.list = res.data.list || []
  myCreateGroupData.count = res.data.count
}

getMyCreateGroup()
getGroup()

function list(){
  getMyCreateGroup()
  getGroup()
}

</script>

<template>
  <div class="contact_view">
    <fim_event event-key="groupList" @event="list"></fim_event>
    <fim_event event-key="friendList" @event="getFriend"></fim_event>
    <div class="contact_slide">
      <fim_slide_head></fim_slide_head>
      <div class="contact_menu">
        <el-scrollbar height="100%">
          <el-menu :default-openeds="['3']">
            <el-sub-menu index="1">
              <template #title>
                <span>我创建的群聊 {{ myCreateGroupData.count }}</span>
              </template>
              <group_list :list="myCreateGroupData.list"></group_list>
            </el-sub-menu>
            <el-sub-menu index="2">
              <template #title>
                <span>我加入的群聊 {{ myJoinGroupData.count }}</span>
              </template>
              <group_list :list="myJoinGroupData.list"></group_list>
            </el-sub-menu>
            <el-sub-menu index="3">
              <template #title>
                <span>我的好友 {{ friendOnlineCount }}/{{ friendData.count }}</span>
              </template>
              <div class="friend_list">
                <div class="item"
                     :class="{active: Number(route.params.id) === item.userID && route.name === 'user_detail'}"
                     @click="checkFriend(item)" v-for="item in friendData.list">
                  <div class="avatar">
                    <img :src="item.avatar" alt="">
                    <div class="online_status" :class="{online: item.isOnline}"></div>
                  </div>
                  <div class="info">
                    <div class="nickname">
                      <el-text style="max-width: 5rem" truncated>{{ item.nickname }}</el-text>
                      （
                      <el-text style="max-width: 4rem" truncated>
                        {{ item.notice === "" ? "-" : item.notice }}
                      </el-text>
                      ）
                    </div>
                    <div class="abstract">
                      <el-text class="w-150px mb-2" truncated>
                        {{ item.abstract }}
                      </el-text>
                    </div>
                  </div>
                </div>
              </div>
            </el-sub-menu>
          </el-menu>
        </el-scrollbar>
      </div>
    </div>
    <div class="contact_main">
      <router-view></router-view>
    </div>
  </div>
</template>

<style lang="scss">
.contact_view {
  width: 100%;
  display: flex;
  height: 100%;

  .contact_slide {
    width: 240px;
    border-right: 1px solid var(--border_color);
    height: 100%;


    .contact_menu {
      height: calc(100% - 40px);
      //overflow-y: auto;

      .el-menu {
        border-right: none;

        .el-sub-menu__title {
          height: 40px;
          font-weight: 600;
          padding: 0 10px;
        }

        .group_list {

        }

        .friend_list {
          width: 100%;

          .item {
            height: 50px;
            display: flex;
            padding: 10px;
            align-items: center;
            cursor: pointer;

            &:hover {
              background-color: var(--item_hover);
            }

            &.active {
              background-color: var(--item_hover);
            }

            .avatar {
              position: relative;
              width: 40px;
              display: flex;
              align-items: center;

              img {
                width: 35px;
                height: 35px;
                border-radius: 5px;
                object-fit: cover;
              }

              .online_status {
                position: absolute;
                right: 6px;
                bottom: 2px;
                width: 5px;
                height: 5px;
                border-radius: 50%;
                background-color: #737373;

                &.online {
                  background-color: #5af50e;
                }
              }
            }

            .info {
              width: calc(100% - 45px);
              font-size: 14px;

              .nickname {
                font-weight: 600;
                display: flex;
                align-items: center;
                color: var(--el-text-color-regular);
              }

              .abstract {
                color: #555;
              }
            }
          }
        }
      }
    }
  }

  .contact_main {
    width: calc(100% - 240px);
  }
}
</style>