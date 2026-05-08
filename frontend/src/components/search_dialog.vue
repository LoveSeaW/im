<script setup lang="ts">
import type {friendAddRequest, userSearchResponse, validResponse} from "@/api/user_api";
import {userValidApi} from "@/api/user_api";
import {reactive} from "vue";
import {ElMessage} from "element-plus";
import {friendAddApi} from "@/api/user_api";
import type {groupAddRequest, groupSearchResponse} from "@/api/group_api";
import {groupAddApi, groupValidApi} from "@/api/group_api";
import  type {baseResponse} from "@/api";
import {ElDialog, ElText, ElInput, ElButton} from "element-plus";
import Fim_avatar from "@/components/fim_avatar.vue";

interface Props {
  visible: boolean
  userInfo?: userSearchResponse
  groupInfo?: groupSearchResponse
  mode: "user" | "group"
}

const props = defineProps<Props>()
const emits = defineEmits(["update:visible", "close"])

function close() {
  emits("update:visible", false)
  emits("close")
}

const data = reactive<validResponse>({
  verification: 0,
  verificationQuestion: {
    problem1: "你童年阴影是什么?",
    problem2: "你小时候玩的好的伙伴叫什么名字？",
    problem3: "你的理想是做什么？",
    answer1: "",
    answer2: "",
    answer3: "",
  }
})

async function getValidData() {
  let res: baseResponse<validResponse> = {code: 0, msg: "", data: data}
  if (props.mode === "user") {
    req.friendId = props.userInfo?.userID as number
    res = await userValidApi(props.userInfo?.userID as number)
  } else {
    reqGroup.groupID = props.groupInfo?.groupId as number
    res = await groupValidApi(props.groupInfo?.groupId as number)
  }
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  Object.assign(data, res.data)
}

const req = reactive<friendAddRequest>({
  friendId: 0,
  verificationQuestion: data.verificationQuestion,
  verify: "",
})
const reqGroup = reactive<groupAddRequest>({
  groupID: 0,
  verificationQuestion: data.verificationQuestion,
  verify: "",
})

async function next() {
  let res: baseResponse<string> = {code: 0, msg: "", data: ""}
  if (props.mode === "user") {
    res = await friendAddApi(req)
  } else {
    res = await groupAddApi(reqGroup)
  }
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  if (props.mode === "user") {
    ElMessage.success("好友申请中...")
  } else {
    ElMessage.success("加群申请中...")
  }
  close()
}


</script>

<template>
  <ElDialog modal-class="search_dialog" @open-auto-focus="getValidData" :model-value="props.visible" @close="close"
             :title="props.mode === 'user' ? '添加好友' : '添加群'" width="500"
             draggable>
    <div class="search_modal_body">
      <div class="left">
        <template v-if="props.mode === 'user'">
          <fim_avatar :size="50" :src="props.userInfo?.avatar as string"></fim_avatar>
        </template>
        <template v-else>
          <fim_avatar :size="50" :src="props.groupInfo?.avatar as string"></fim_avatar>
        </template>
        <ElText truncated style="max-width: 6rem;">
          {{ props.mode === 'user' ? props.userInfo?.nickname : props.groupInfo?.title }}
        </ElText>
        <div class="group_info" v-if="props.mode === 'group'">
          <div>成员：{{ props.groupInfo?.userCount }}人</div>
          <div>简介：{{ props.groupInfo?.abstract }}</div>
        </div>
      </div>
      <div class="right">

        <div class="mode" v-if="data.verification === 2">
          <div class="label">请输入验证消息</div>
          <ElInput v-if="props.mode === 'user'" type="textarea" resize="none" v-model="req.verify" :autosize="{minRows: 3, maxRows: 4}"
                    placeholder="我是..."></ElInput>
          <ElInput v-else type="textarea" resize="none" v-model="reqGroup.verify" :autosize="{minRows: 3, maxRows: 4}"
                    placeholder="我是..."></ElInput>
        </div>
        <div class="mode" v-else-if="data.verification === 0">
          <div class="label">{{ props.mode === "user" ? "该用户不允许被加好友" : "不允许任何人加群" }}</div>
        </div>
        <div class="mode" v-else-if="data.verification === 1">
          <div class="label">{{ props.mode === "user" ? "可直接加好友" : "可直接加群" }}</div>
        </div>
        <div class="mode" v-else-if="data.verification === 3 || data.verification === 4">
          <div class="label">
            {{ data.verification === 3 ? "请回答验证问题" : "请正确回答验证问题" }}
          </div>
          <div class="item" v-if="data.verificationQuestion.problem1">
            <span class="label">问题1</span>
            <span class="question">{{ data.verificationQuestion.problem1 }}</span>
            <ElInput placeholder="答案" v-model="data.verificationQuestion.answer1"></ElInput>
          </div>
          <div class="item" v-if="data.verificationQuestion.problem2">
            <span class="label">问题2</span>
            <span class="question">{{ data.verificationQuestion.problem2 }}</span>
            <ElInput placeholder="答案" v-model="data.verificationQuestion.answer2"></ElInput>
          </div>
          <div class="item" v-if="data.verificationQuestion.problem3">
            <span class="label">问题3</span>
            <span class="question">{{ data.verificationQuestion.problem3 }}</span>
            <ElInput placeholder="答案" v-model="data.verificationQuestion.answer3"></ElInput>
          </div>
        </div>
      </div>
    </div>
    <template #footer>
      <ElButton size="small" type="primary" @click="next()" :disabled="data.verification === 0">下一步</ElButton>
      <ElButton size="small" @click="close">取消</ElButton>
    </template>
  </ElDialog>
</template>

<style lang="scss">
.search_dialog {
  .el-dialog {
    padding: 0;

    .el-dialog__header {
      padding: 10px 20px;
      background-color: var(--active);

      .el-dialog__title {
        color: white;
        font-size: 14px;
      }

      .el-dialog__headerbtn {
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;

        i {
          color: white;
          font-weight: 600;
        }
      }
    }

    .search_modal_body {
      height: 300px;
      width: 100%;
      display: flex;

      .left {
        width: 120px;
        background-color: #f5f5f5;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        padding-top: 20px;

        .el-text {
          margin-top: 10px;
        }

        .group_info {
          padding: 10px;
        }
      }

      .right {
        width: calc(100% - 120px);
        padding: 20px;

        .mode > .label {
          margin-bottom: 5px;
        }

        .item {
          .label {
            &::after {
              content: ":";
              display: inline-block;
              margin-right: 5px;
            }
          }

          .el-input {
            display: flex;
            margin-top: 5px;
            margin-bottom: 10px;
          }
        }
      }
    }

    .el-dialog__footer {
      background-color: #e6e6e6;
      display: flex;
      align-items: center;
      padding-top: 0;
      height: 40px;
      justify-content: end;
      padding-right: 20px;
    }
  }
}
</style>