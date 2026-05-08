<script setup lang="ts">
import {useRoute} from "vue-router";
import {useStore} from "@/stores";
import {reactive, ref} from "vue";
import {registerApi, type registerRequest} from "@/api/auth_api";
import {ElMessage, type FormRules} from "element-plus";
import router from "@/router";

const route = useRoute()
const store = useStore()
const form = reactive<registerRequest>({
  nickname: "", pwd: "", rePwd: ""
});

const rules = reactive<FormRules>({
  nickname: [
    {required: true, message: '请输入昵称', trigger: 'blur'},
  ],
  pwd: [
    {required: true, message: '请输入密码', trigger: 'blur'},
  ],
  rePwd: [
    {required: true, message: '请输入确认密码', trigger: 'blur'},
  ]
})

const formRef = ref()

async function register() {
  let val = await formRef.value.validate()
  if (!val) {
    return
  }

  let res = await registerApi(form)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  // 拿到的是token，前端要对他进行解码
  ElMessage.success("注册成功")
  router.push({
    name: "registerSuccess",
    query: {
      userID: res.data.userID,
    }
  })
}
</script>

<template>
  <el-form ref="formRef" :model="form" :rules="rules">
    <el-form-item prop="nickname">
      <el-input v-model="form.nickname" placeholder="昵称">
        <template #prefix>
          <i class="iconfont icon-yonghuming"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="pwd">
      <el-input v-model="form.pwd" type="password" placeholder="密码">
        <template #prefix>
          <i class="iconfont icon-mima"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="rePwd">
      <el-input v-model="form.rePwd" type="password" placeholder="确认密码">
        <template #prefix>
          <i class="iconfont icon-mima"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="item_action xxjfgdfgg">
      <span></span>
      <router-link :to="{name: 'login'}">已有账号，去登录</router-link>
    </el-form-item>
    <el-form-item class="item_btn">
      <el-button style="width: 100%" @click="register" type="primary">申请</el-button>
    </el-form-item>
  </el-form>
</template>

<style lang="scss">
.xxjfgdfgg {
  .el-form-item__content {
    display: flex;
    justify-content: space-between;
  }

}
</style>