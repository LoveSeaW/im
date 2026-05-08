<script setup lang="ts">
import {useStore} from "@/stores";
import {ElMessage} from "element-plus";
import {reactive, watch} from "vue";
import type {verificationQuestionType} from "@/api/user_api";
import {type userConfUpdateType, userInfoUpdateApi} from "@/api/user_api";

const store = useStore()
// 0 不允许别人查找到我， 1  通过用户号找到我 2 可以通过昵称搜索到我
// 0 不允许任何人添加  1 允许任何人添加  2 需要验证消息 3 需要回答问题  4  需要正确回答问题

const verificationQuestion = reactive<verificationQuestionType>({
  problem1: "",
  problem2: "",
  problem3: "",
  answer1: "",
  answer2: "",
  answer3: "",
})

async function save() {
  let res = await userInfoUpdateApi({
    searchUser: store.userConfInfo.searchUser,
    verification: store.userConfInfo.verification,
    verificationQuestion: verificationQuestion,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("防骚扰设置修改成功")
}

watch(() => store.userConfInfo.verificationQuestion, () => {
  Object.assign(verificationQuestion, store.userConfInfo.verificationQuestion)
}, {deep: true})


</script>

<template>
  <div class="role_info_view">
    <el-form-item label="防骚扰">
      <div>
        <div>
          <div class="label">请选择允许别人查找到你的方式：</div>
          <el-radio-group v-model="store.userConfInfo.searchUser">
            <el-radio :value="0">不允许别人查找到我</el-radio>
            <el-radio :value="1">只能通过用户号查找到我</el-radio>
            <el-radio :value="2">可以通过昵称搜索到我</el-radio>
          </el-radio-group>
        </div>
        <div>
          <div class="label">请选择适合你的验证方式：</div>
          <el-radio-group v-model="store.userConfInfo.verification">
            <el-radio :value="0">不允许任何人添加</el-radio>
            <el-radio :value="1">允许任何人添加</el-radio>
            <el-radio :value="2">需要验证消息</el-radio>
            <el-radio :value="3">需要回答问题</el-radio>
            <el-radio :value="4">需要正确回答问题</el-radio>
          </el-radio-group>
          <div class="verification"
               v-if="store.userConfInfo.verification === 3 || store.userConfInfo.verification === 4">
            <div class="item">
              <div class="question"><span>问题1：</span>
                <el-input placeholder="问题1" v-model="verificationQuestion.problem1"></el-input>
              </div>
              <div class="answer" v-if="store.userConfInfo.verification === 4"><span>答案:</span>
                <el-input placeholder="答案" v-model="verificationQuestion.answer1"></el-input>
              </div>
            </div>
            <div class="item">
              <div class="question"><span>问题2：</span>
                <el-input placeholder="问题2" v-model="verificationQuestion.problem2"></el-input>
              </div>
              <div class="answer" v-if="store.userConfInfo.verification === 4"><span>答案:</span>
                <el-input placeholder="答案" v-model="verificationQuestion.answer2"></el-input>
              </div>
            </div>
            <div class="item">
              <div class="question"><span>问题3：</span>
                <el-input placeholder="问题2" v-model="verificationQuestion.problem3"></el-input>
              </div>
              <div class="answer" v-if="store.userConfInfo.verification === 4"><span>答案:</span>
                <el-input placeholder="答案" v-model="verificationQuestion.answer3"></el-input>
              </div>
            </div>
          </div>

        </div>
      </div>
    </el-form-item>
    <el-form-item label="保存">
      <el-button type="primary" @click="save">保存</el-button>
    </el-form-item>
  </div>
</template>

<style lang="scss">
.role_info_view {
  .el-radio-group {
    display: flex;
    flex-direction: column;
    align-items: start;
  }

  .question, .answer {
    display: flex;
    width: 500px;
    margin-bottom: 10px;

    .el-input {
      width: 300px;
    }

    > span {
      width: 4rem;
    }
  }

}
</style>