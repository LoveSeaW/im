<script setup lang="ts">
import {useStore} from "@/stores";
import {ref, nextTick, watch} from "vue";
import {type userConfUpdateType, userInfoUpdateApi} from "@/api/user_api";
import {ElMessage} from "element-plus";
import Avatar_cropper from "@/components/avatar_cropper.vue";

export interface IProps {
  type: string // 上传类型, 企业logo / 浏览器logo
  allowTypeList: string[] // 接收允许上传的图片类型
  limitSize: number // 限制大小
  fixedNumber: number[] // 截图框的宽高比例
  fixedNumberAider?: number[] // 侧边栏收起截图框的宽高比例
  previewWidth: number // 预览宽度
  title?: string // 裁剪标题
}

const store = useStore()

interface iptType {
  label: string
  isShowIpt: boolean
  maxlength?: number
  rows?: number
  type: "text" | "textarea" | "password" | "button" | "checkbox" | "file" | "number" | "radio"
  val: string
  old: string
  key: "nickname" | "abstract"
}

const list = ref<iptType[]>([
  {
    label: "昵称",
    isShowIpt: false,
    maxlength: 13,
    val: store.userConfInfo.nickname,
    type: "text",
    old: "",
    key: "nickname"
  },
  {
    label: "简介",
    isShowIpt: false,
    type: "textarea",
    rows: 3,
    old: "",
    val: store.userConfInfo.abstract,
    key: "abstract"
  }
])
const editRefList = ref()

watch(() => store.userConfInfo, () => {
  list.value[0].val = store.userConfInfo.nickname
  list.value[1].val = store.userConfInfo.abstract
}, {deep: true})

function edit(index: number) {
  list.value[index].isShowIpt = true
  list.value[index].old = list.value[index].val
  nextTick(() => {
    if (editRefList.value.length) {
      editRefList.value[0].focus()
    }
  })
}

async function blur(index: number) {
  list.value[index].isShowIpt = false
  if (list.value[index].old == list.value[index].val) {
    return
  }
  let data: userConfUpdateType = {}

  data[list.value[index].key] = list.value[index].val
  let res = await userInfoUpdateApi(data)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success(list.value[index].label + "修改成功")
}

const clipperData = ref<IProps>({
  type: '',
  allowTypeList: [],
  limitSize: 1,
  fixedNumber: [],
  previewWidth: 0
})
const clipperRef = ref()

function showCropper() {
  clipperData.value = {
    type: 'browserLogo', // 该参数可根据实际要求修改类型
    allowTypeList: ['png', 'jpg', 'jpeg'], // 允许上传的图片格式
    limitSize: 1, // 限制的大小
    fixedNumber: [1, 1],  // 截图比例，可根据实际情况进行修改
    previewWidth: 100, // 预览宽度
  }
  // 打开裁剪组件
  clipperRef.value.uploadFile()
}

async function onConfirm(val: any) {

  let res = await userInfoUpdateApi({
    avatar: val
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  // 马上把用户信息的头像变掉
  store.userConfInfo.avatar = val
  store.userInfo.avatar = val
  store.saveToken()
  ElMessage.success("用户头像更新成功")
}


</script>

<template>
  <div class="my_info_view">
    <avatar_cropper ref="clipperRef"
                    :type="clipperData.type"
                    :allow-type-list="clipperData.allowTypeList"
                    :limit-size="clipperData.limitSize"
                    :fixed-number="clipperData.fixedNumber"
                    :preview-width="clipperData.previewWidth"
                    @confirm="onConfirm"></avatar_cropper>
    <el-form-item label="头像">
      <el-avatar :src="store.userConfInfo.avatar" @click="showCropper"></el-avatar>
    </el-form-item>
    <el-form-item label="用户号">
      <span>{{ store.userConfInfo.userID }}</span>
    </el-form-item>
    <el-form-item :label="item.label" v-for="(item, index) in list">
      <span v-if="!item.isShowIpt">{{ item.val }}</span>
      <el-input v-else ref="editRefList" :maxlength="item.maxlength as number" :rows="item.rows" :type="item.type"
                @blur="blur(index)" class="edit_ipt"
                v-model="item.val" :placeholder="'修改'+item.label"></el-input>
      <i class="iconfont icon-bianji" @click="edit(index)"></i>
    </el-form-item>
  </div>
</template>

<style>
.my_info_view {
  i {
    cursor: pointer;
    margin-left: 5px;
  }

  .edit_ipt {
    display: inline-flex;
    width: 200px;
  }
}
</style>