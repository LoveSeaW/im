<script setup lang="ts">
import Fim_avatar from "@/components/fim_avatar.vue";
import {useStore} from "@/stores";
import {Edit} from "@element-plus/icons-vue";
import {nextTick, ref} from "vue";
import {userInfoUpdateApi} from "@/api/user_api";
import {ElMessage} from "element-plus";
import type {IProps} from "@/views/web/info/my_info.vue";
import Avatar_cropper from "@/components/avatar_cropper.vue";

const store = useStore()

const isEditName = ref(false)
const editNameRef = ref()
const oldName = ref("")

function showEditName(old: string) {
  isEditName.value = true
  oldName.value = old
  nextTick(() => {
    editNameRef.value.focus()
  })
}

function editNameBlur() {
  isEditName.value = false
  if (oldName.value != store.groupData.title) {
    store.updateGroup()
  }
}

const isEditAbstract = ref(false)
const editAbstractRef = ref()
const oldAbstract = ref("")

function showEditAbstract(old: string) {
  isEditAbstract.value = true
  oldAbstract.value = old
  nextTick(() => {
    editAbstractRef.value.focus()
  })
}

function editAbstractBlur() {
  isEditAbstract.value = false
  if (oldAbstract.value != store.groupData.abstract) {
    store.updateGroup()
  }
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
  if (store.groupData.role === 3) {
    return
  }
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
  store.groupData.avatar = val
  store.updateGroup()
}


</script>

<template>
  <div class="information_view">
    <el-form>
      <avatar_cropper ref="clipperRef"
                      image-type="group_avatar"
                      :type="clipperData.type"
                      :allow-type-list="clipperData.allowTypeList"
                      :limit-size="clipperData.limitSize"
                      :fixed-number="clipperData.fixedNumber"
                      :preview-width="clipperData.previewWidth"
                      @confirm="onConfirm"></avatar_cropper>
      <el-form-item label="群头像">
        <fim_avatar :src="store.groupData.avatar" @click="showCropper"></fim_avatar>
      </el-form-item>
      <el-form-item label="群名称">
        <el-input style="width: 200px" ref="editNameRef" @blur="editNameBlur" v-if="isEditName" placeholder="修改群名称"
                  v-model="store.groupData.title"></el-input>
        <span v-else>{{ store.groupData.title }}</span>
        <el-icon :size="18" class="icon" @click="showEditName(store.groupData.title)" v-if="store.groupData.role != 3">
          <Edit></Edit>
        </el-icon>
      </el-form-item>
      <el-form-item label="群介绍">
        <el-input type="textarea" style="width: 400px" :autosize="{minRows: 3, maxRows: 4}" ref="editAbstractRef"
                  @blur="editAbstractBlur" v-if="isEditAbstract" placeholder="修改群名称"
                  v-model="store.groupData.abstract"></el-input>
        <el-text v-else truncated>{{ store.groupData.abstract }}</el-text>
        <el-icon :size="18" class="icon" @click="showEditAbstract(store.groupData.abstract)"
                 v-if="store.groupData.role != 3">
          <Edit></Edit>
        </el-icon>
      </el-form-item>
      <el-form-item label="群主/管理员"></el-form-item>
      <div class="admins">
        <div class="creater">
          <fim_avatar :src="store.groupData.creator.avatart" :title="store.groupData.creator.nickname"></fim_avatar>
        </div>
        <div class="admin_list">
          <div class="admin" v-for="item in store.groupData.adminList">
            <fim_avatar :src="item.avatart" :title="item.nickname"></fim_avatar>
          </div>
        </div>
      </div>
    </el-form>
  </div>
</template>

<style lang="scss">
.information_view {
  .el-form-item {
    margin-bottom: 5px;

    .el-form-item__label {
      font-weight: 600;
    }

    .icon {
      cursor: pointer;
      color: var(--text_color);
      margin-left: 5px;
    }
  }

  .admins {
    display: flex;
    justify-content: start;

    .creater {
      border-right: 1px solid var(--border_color);
      padding-right: 10px;
    }

    .admin {
      margin-right: 10px;
    }
  }

}
</style>