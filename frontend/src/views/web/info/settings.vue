<script setup lang="ts">
import {useStore} from "@/stores";

const store = useStore()

function themeChange(){
  console.log(store.theme)
  store.setTheme(store.theme)
}

</script>

<template>
  <div class="my_info_view">
    <el-form-item label="主题设置">
      <el-checkbox @change="themeChange" v-model="store.theme">{{ store.theme ? "白天模式" : "暗黑模式" }}</el-checkbox>
    </el-form-item>
    <el-form-item label="退出登陆">
      <el-button type="danger" size="small" @click="store.logout">退出登陆</el-button>
    </el-form-item>
  </div>
</template>