<script setup lang="ts">
import {dateTimeFormat} from "@/utils/date";
import type {verificationQuestionType} from "@/api/user_api";
import {DArrowRight} from "@element-plus/icons-vue";

export interface ValidType {
  sendAvatar: string
  sendName: string
  revAvatar: string
  revName: string
  createdAt: string
  additionalMessages: string
  verificationQuestion?: verificationQuestionType
}

interface Props {
  data: ValidType
}

const props = defineProps<Props>()


</script>

<template>
  <div class="notice_detail">
    <div class="head">
      <div class="send">
          <el-avatar :size="50" :src="props.data.sendAvatar"></el-avatar>
          <span>{{ props.data.sendName }}</span>
      </div>
      <div class="center">
        <el-icon size="20">
          <DArrowRight/>
        </el-icon>
        <span>申请</span>
      </div>
      <div class="rev">
          <el-avatar :size="50" v-if="props.data.revAvatar.startsWith('/api') || props.data.revAvatar.startsWith('http')" :src="props.data.revAvatar"></el-avatar>
          <el-avatar :size="50" v-else style="background-color: #529b2e">{{ props.data.revAvatar }}</el-avatar>
          <span>{{ props.data.revName }}</span>
      </div>
    </div>
    <div class="body">
      <div class="item">
        <span class="label">申请时间</span>
        <span class="val">{{ dateTimeFormat(props.data.createdAt) }}</span>
      </div>
      <div class="item">
        <span class="label">附加消息</span>
        <span class="val">{{ props.data.additionalMessages }}</span>
      </div>
      <div class="item">
        <span class="label">验证消息</span>
      </div>
      <div class="verificationQuestion" v-if="props.data.verificationQuestion">
        <div class="item" v-if="props.data.verificationQuestion.problem1">
          <span class="label">问题1</span>
          <span class="question">{{ props.data.verificationQuestion?.problem1 }}</span>
          <div>
            <span class="label">答案</span>
            <span class="answer">{{ props.data.verificationQuestion?.answer1 }}</span>
          </div>
        </div>
        <div class="item" v-if="props.data.verificationQuestion.problem2">
          <span class="label">问题2</span>
          <span class="question">{{ props.data.verificationQuestion?.problem2 }}</span>
          <div>
            <span class="label">答案</span>
            <span class="answer">{{ props.data.verificationQuestion?.answer2 }}</span>
          </div>

        </div>
        <div class="item" v-if="props.data.verificationQuestion.problem3">
          <span class="label">问题3</span>
          <span class="question">{{ props.data.verificationQuestion?.problem2 }}</span>
          <div>
            <span class="label">答案</span>
            <span class="answer">{{ props.data.verificationQuestion?.answer3 }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.notice_detail {
  padding: 20px;

  .head {
    background-color: #f9f9f9;
    display: flex;
    justify-content: center;
    height: 100px;
    align-items: center;
    font-size: 14px;
    color: var(--text_color);

    > div {
      display: flex;
      flex-direction: column;
      align-items: center;
    }

    .center {
      width: 100px;
    }
  }

  .body {
    margin-top: 20px;

    .item {
      margin-bottom: 10px;
      color: var(--text_color);

      .label {
        &::after {
          content: ":";
          display: inline-block;
          margin-right: 5px;
        }
      }
    }
  }
}
</style>