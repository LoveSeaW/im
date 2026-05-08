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
    {required: true, message: '请输入管理员账号', trigger: 'blur'},
  ],
  password: [
    {required: true, message: '请输入管理员密码', trigger: 'blur'},
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
  
  ElMessage.success("管理员登录成功")
  store.setToken(res.data.token)

  router.push({
    name: "admin_dashboard",
  })
}
</script>

<template>
  <el-form ref="formRef" :model="form" :rules="rules">
    <el-form-item prop="userName">
      <el-input v-model="form.userName" placeholder="管理员账号">
        <template #prefix>
          <i class="iconfont icon-yonghuming"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password" class="item_password">
      <el-input v-model="form.password" type="password" placeholder="管理员密码">
        <template #prefix>
          <i class="iconfont icon-mima"></i>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item class="item_action xxjfgdfgg">
      <el-checkbox>记住密码</el-checkbox>
    </el-form-item>
    <el-form-item class="item_btn">
      <el-button style="width: 100%" @click="login" type="primary">管理员登录</el-button>
    </el-form-item>
  </el-form>
</template>

<style lang="scss">
.xxjfgdfgg{
  .el-form-item__content{
    display: flex;
    justify-content: flex-start;
  }
}
</style>