<script setup lang="ts">
import {useRoute, useRouter} from "vue-router";

const route = useRoute()
const router = useRouter()


function toLogin() {
  router.push({name: "login", query: {userID: route.query.userID}})
}

</script>

<template>
  <el-form ref="formRef" class="register_success_form">
    <div>
      <div>
        <i class="iconfont icon-chenggong"></i>
      </div>
      <div>
        <div>申请成功，您获得用户号：</div>
        <div>{{ route.query.userID }}</div>
      </div>
    </div>
    <el-form-item class="item_btn">
      <el-button style="width: 100%" @click="toLogin" type="primary">去登录</el-button>
    </el-form-item>
  </el-form>
</template>

<style lang="scss">
.register_success_form {
  > div {
    display: flex;
    padding: 20px 0;

    > div:nth-child(1) {
      width: 50px;

      i {
        font-size: 50px;
        color: #1fc41f;
      }
    }

    > div:nth-child(2) {
      margin-left: 20px;

      div:nth-child(1) {
        font-size: 14px;
        color: var(--text_color);
      }

      > div:nth-child(2) {
        font-size: 30px;
        color: var(--text_color);
      }
    }
  }
}
</style>