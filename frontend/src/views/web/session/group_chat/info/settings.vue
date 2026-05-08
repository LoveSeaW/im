<script setup lang="ts">
import {reactive, watch} from "vue";
import type {groupUpdateRequest} from "@/api/group_api";
import {useStore} from "@/stores";
import {groupUpdateInfoApi} from "@/api/group_api";
import {ElMessage} from "element-plus";
import type {verificationQuestionType} from "@/api/user_api";
import {groupRemoveApi} from "@/api/group_api";
import router from "@/router";

const store = useStore()
const data = reactive<groupUpdateRequest>({
  id: store.groupData.groupId,
  isSearch: false,
  verification: 2,
  isInvite: false,
  isTemporarySession: false,
  verificationQuestion: {
    problem1: "",
    problem2: "",
    problem3: "",
    answer1: "",
    answer2: "",
    answer3: ""
  }
})

watch(() => store.groupData, () => {
  data.id = store.groupData.groupId
  data.isProhibition = store.groupData.isProhibition
  data.isSearch = store.groupData.isSearch
  data.isInvite = store.groupData.isInvite
  data.isTemporarySession = store.groupData.isTemporarySession
}, {immediate: true, deep: true})

async function updateGroup() {
  let res = await groupUpdateInfoApi(data)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("修改群资料成功")
}

async function dissolution() {
  let res = await groupRemoveApi(data.id)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("群解散成功")
  router.push({
    name: "session"
  })
}

</script>

<template>
  <div class="group_settings">
    <el-scrollbar height="500px">
      <el-form :model="data">
        <el-form-item label="查找方式">
          <div>
            <div>
              <el-checkbox v-model="data.isSearch">允许被搜索</el-checkbox>
            </div>
            <div style="font-size: 14px; color: var(--text_color)">
              <span>（勾选后将启用群号，他人可按群号、群名称搜索本群）</span>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="加群方式">
          <div>
            <div>
              <el-radio-group class="radio_group" v-model="data.verification">
                <el-radio :value="0">不允许任何人加群</el-radio>
                <el-radio :value="1">允许任何人加群</el-radio>
                <el-radio :value="2">需要验证消息</el-radio>
                <el-radio :value="3">需要回答问题</el-radio>
                <el-radio :value="4">需要正确回答问题</el-radio>
              </el-radio-group>
            </div>
            <div class="verificationQuestion" v-if="data.verification === 3 || data.verification === 4">
              <div class="item">
                <div>
                  <span class="label">问题1</span>
                  <el-input placeholder="问题1"
                            v-model="(data.verificationQuestion as verificationQuestionType).problem1"></el-input>
                </div>
                <div v-if="data.verification === 4">
                  <span class="label">答案</span>
                  <el-input placeholder="答案"
                            v-model="(data.verificationQuestion as verificationQuestionType).answer1"></el-input>
                </div>
              </div>
              <div class="item">
                <div>
                  <span class="label">问题2</span>
                  <el-input placeholder="问题2"
                            v-model="(data.verificationQuestion as verificationQuestionType).problem2"></el-input>
                </div>
                <div v-if="data.verification === 4">
                  <span class="label">答案</span>
                  <el-input placeholder="答案"
                            v-model="(data.verificationQuestion as verificationQuestionType).answer2"></el-input>
                </div>
              </div>
              <div class="item">
                <div>
                  <span class="label">问题3</span>
                  <el-input placeholder="问题3"
                            v-model="(data.verificationQuestion as verificationQuestionType).problem3"></el-input>
                </div>
                <div v-if="data.verification === 4">
                  <span class="label">答案</span>
                  <el-input placeholder="答案"
                            v-model="(data.verificationQuestion as verificationQuestionType).answer3"></el-input>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="邀请方式">
          <el-checkbox v-model="data.isInvite">允许群成员邀请好友进群</el-checkbox>
        </el-form-item>
        <el-form-item label="会话权限">
          <div class="radio_group">
            <el-checkbox v-model="data.isTemporarySession" disabled>允许群成员发起临时会话</el-checkbox>
            <el-checkbox v-model="data.isProhibition">全员禁言</el-checkbox>
          </div>

        </el-form-item>
        <el-form-item label="其他">
          <el-button size="small" type="primary" @click="updateGroup">修改群设置</el-button>
          <el-popconfirm v-if="store.groupData.role === 1" title="是否解散该群？"
                         @confirm="dissolution">
            <template #reference>
              <el-button size="small" type="danger">解散该群</el-button>
            </template>
          </el-popconfirm>
        </el-form-item>
      </el-form>
    </el-scrollbar>
  </div>
</template>

<style lang="scss">
.group_settings {
  .el-form-item {
    margin-bottom: 8px;
  }

  .radio_group {
    display: flex;
    flex-direction: column;
    align-items: start;
  }

  .verificationQuestion {
    .item {
      margin-bottom: 10px;

      > div {
        display: flex;
        margin-bottom: 5px;

        .label {
          width: 3rem;
          flex-shrink: 0;
          margin-right: 5px;
        }
      }
    }
  }
}
</style>