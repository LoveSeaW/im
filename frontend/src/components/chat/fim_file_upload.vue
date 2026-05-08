<script setup lang="ts">

import {reactive, ref} from "vue";
import {fileSizeFormat} from "@/utils/file";
import {uploadFileApi} from "@/api/file_api";
import {ElMessage} from "element-plus";
import {useStore} from "@/stores";
import Fim_file_icon from "@/components/chat/fim_file_icon.vue";
import {Msg} from "@/service/send_msg";

const store = useStore()

interface Props {
  type: "user" | "group"
  id: number // 用户就是用户id，群就是群id
}

const props = defineProps<Props>()
const msg = new Msg(props.type)
interface selectFileType {
  fileName: string
  fileSize: number
  src: string
  file?: File
}

const selectImg = reactive<selectFileType>({
  fileName: "",
  fileSize: 0,
  src: "",
  file: undefined,
})

const imgVisible = ref(false)


const iptRef = ref()

function fileClick() {
  (iptRef.value as HTMLInputElement).click();
}

// 选中了文件
function fileIptChange(e: Event) {
  if (((e.target as HTMLInputElement).files as FileList).length > 0) {
    const file = ((e.target as HTMLInputElement).files as FileList)[0]

    selectImg.fileName = file.name
    selectImg.file = file
    selectImg.fileSize = file.size
    imgVisible.value = true
  }
}

// 发送图片
async function sendImage() {
  let res = await uploadFileApi(selectImg.file as File)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  msg.fileMsg(props.id, res.data.src)
  cancel()
}

function cancel() {
  selectImg.src = ""
  selectImg.fileName = ""
  selectImg.fileSize = 0
  imgVisible.value = false
}




</script>

<template>
  <div class="fim_file_upload">
    <el-dialog width="30%" title="文件上传" v-model="imgVisible" modal-class="file_dialog">
      <div class="image_preview_body">
        <fim_file_icon :name="selectImg.fileName"></fim_file_icon>
        <div class="info">
          <div class="name">{{ selectImg.fileName }}</div>
          <div class="size">{{ fileSizeFormat(selectImg.fileSize, 2) }}</div>
        </div>
        <i class="iconfont icon-quxiao" @click="cancel"></i>
      </div>
      <template #footer>
        <el-button @click="sendImage" type="primary">发送</el-button>
        <el-button @click="cancel">取消</el-button>
      </template>
    </el-dialog>
    <input type="file" ref="iptRef" @change="fileIptChange" style="display: none">
    <slot><i title="文件上传" class="iconfont icon-wenjian" @click="fileClick"></i></slot>
  </div>

</template>

<style lang="scss">
.file_dialog {
  .image_preview_body {
    display: flex;
    align-items: center;
    background-color: var(--item_hover);
    position: relative;
    padding: 20px;
    border-radius: 5px;

    .iconfont{
      font-size: 24px;
    }


    .info{
      margin-left: 10px;
    }
    &:hover {
      .icon-quxiao {
        opacity: 1;
      }
    }

    .icon-quxiao {
      position: absolute;
      right: 10px;
      opacity: 0;
      cursor: pointer;
    }
  }
}
</style>