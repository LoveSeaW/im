<script setup lang="ts">
import {reactive, ref} from "vue";
import {groupCreateApi} from "@/api/group_api";
import type {groupCreateRequest} from "@/api/group_api";
import {ElMessage, type FormRules} from "element-plus";
import {friendListApi, type friendType} from "@/api/user_api";
import {eventStore} from "@/stores/event_store";

const store = eventStore()
const data = reactive<groupCreateRequest>({
  "mode": 1,
  "name": "",
  "size": 10,
  isSearch: false,
  userIdList: []
})
const rules = reactive<FormRules<string>>({
  name: [
    {required: true, message: '请输入群名称', trigger: 'blur'},
  ],
})

const myFriendList = ref<friendType[]>([])

async function getData() {
  const res = await friendListApi({limit: -1})
  myFriendList.value = res.data.list || []
}

getData()

const formRef = ref()

async function groupCreate() {

  if (data.mode === 1) {
    const val = await formRef.value.validate()
    if (val) {
      return
    }

  }

  let res = await groupCreateApi(data)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("创建成功")
  store.setEventKey("groupList")

}

</script>

<template>
  <div class="add_group_view">
    <div class="head">
      <span @click="data.mode = 1" :class="{active:data.mode === 1}">直接创建</span>
      <span @click="data.mode = 2" :class="{active:data.mode === 2}">选人创建</span>
    </div>
    <div class="body">
      <div v-if="data.mode === 1">
        <el-form ref="formRef" :model="data" :rules="rules">
          <el-form-item label="群名称" prop="name">
            <el-input placeholder="请输入群名称" v-model="data.name"></el-input>
          </el-form-item>
          <el-form-item label="群搜索">
            <el-checkbox v-model="data.isSearch">可以通过群名称搜索到</el-checkbox>
          </el-form-item>
          <el-form-item label="群规模">
            <el-radio-group v-model="data.size">
              <el-radio :value="10">10人群</el-radio>
              <el-radio :value="100">100人群</el-radio>
              <el-radio :value="200">200人群</el-radio>
            </el-radio-group>
          </el-form-item>
          <el-button type="primary" @click="groupCreate">创建群</el-button>
        </el-form>
      </div>
      <div v-if="data.mode === 2">
        <el-checkbox-group v-model="data.userIdList">
          <el-checkbox :value="item.userID" v-for="item in myFriendList">
            <el-avatar :src="item.avatar" :size="40"></el-avatar>
            <span class="nickname">{{ item.nickname }}</span>
          </el-checkbox>
        </el-checkbox-group>
        <el-button type="primary" @click="groupCreate">创建群</el-button>
      </div>
    </div>
  </div>

</template>

<style lang="scss">
.add_group_view {
  padding: 10px;

  .head {
    padding-bottom: 10px;
    border-bottom: 1px solid var(--border_color);

    span {
      margin-right: 20px;
      font-size: 16px;
      color: var(--text_color);
      cursor: pointer;

      &.active {
        color: var(--active);
      }
    }
  }

  .body {
    padding-top: 10px;
    width: 70%;

    .el-form-item {
      margin-bottom: 10px;
    }
    .el-checkbox-group{
      display: flex;
      flex-direction: column;
      .el-checkbox{
        height: 50px;
        .el-checkbox__label{
          display: flex;
          align-items: center;
          .nickname{
            margin-left: 10px;
          }
        }
      }
    }
  }
}
</style>