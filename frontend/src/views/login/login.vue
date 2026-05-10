<script setup lang="ts">
import {useRoute} from "vue-router";
import {useStore} from "@/stores";
import {reactive, ref} from "vue";
import {authLoginApi, type authLoginRequest} from "@/api/auth_api";
import {ElMessage, type FormRules} from "element-plus";
import router from "@/router";

const route = useRoute()
const store = useStore()
const form = reactive<authLoginRequest>({
  userName: route.query.userID as string | "",
  password: "",
});

const rules = reactive<FormRules>({
  userName: [
    {required: true, message: '请输入账号/用户ID', trigger: 'blur'},
  ],
  password: [
    {required: true, message: '请输入密码', trigger: 'blur'},
  ]
})

const formRef = ref()

async function login() {
  let val = await formRef.value.validate()
  if (!val) {
    return
  }

  let res = await authLoginApi(form)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  
  ElMessage.success("登录成功")
  store.setToken(res.data.token)

  router.push({
    name: "web",
  })
}
</script>

<template>
  <el-form ref="formRef" :model="form" :rules="rules">
    <el-form-item prop="userName">
      <el-input v-model="form.userName" placeholder="账号/用户ID">
        <template #prefix>
          <i class="iconfont icon-yonghuming"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password" class="item_password">
      <el-input v-model="form.password" type="password" placeholder="密码">
        <template #prefix>
          <i class="iconfont icon-mima"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="item_action xxjfgdfgg">
      <span></span>
      <router-link :to="{name: 'register'}">没有账号？去注册</router-link>
    </el-form-item>
    <el-form-item class="item_btn">
      <el-button style="width: 100%" @click="login" type="primary">登录</el-button>
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