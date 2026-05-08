<script setup lang="ts">

import {reactive, ref} from "vue";
import {fileSizeFormat} from "@/utils/file";
import {uploadImageApi} from "@/api/file_api";
import {ElMessage} from "element-plus";
import {useStore} from "@/stores";
import {Msg} from "@/service/send_msg";

const store = useStore()

interface Props {
  type: "user" | "group"
  id: number // 用户就是用户id，群就是群id
}

const props = defineProps<Props>()
const msg = new Msg(props.type)
interface selectImgType {
  fileName: string
  fileSize: number
  src: string
  file?: File
}

const selectImg = reactive<selectImgType>({
  fileName: "",
  fileSize: 0,
  src: "",
  file: undefined,
})

const imgVisible = ref(false)


const imgIptRef = ref()

function imgClick() {
  (imgIptRef.value as HTMLInputElement).click();
}

// 选中了图片
function imageIptChange(e: Event) {
  if (((e.target as HTMLInputElement).files as FileList).length > 0) {
    const file = ((e.target as HTMLInputElement).files as FileList)[0]
    let reader = new FileReader();
    reader.onload = (e: ProgressEvent<FileReader>) => {
      let base64 = (e.target as FileReader).result as string;
      selectImg.src = base64
      selectImg.fileName = file.name
      selectImg.file = file
      selectImg.fileSize = file.size
      imgVisible.value = true
    }
    reader.readAsDataURL(file)
  }
}

// 发送图片
async function sendImage() {
  let res = await uploadImageApi(selectImg.file as File, "chat")
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  msg.imageMsg(props.id,{
    title: selectImg.fileName,
    src:  res.data.url,
  })
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
  <div class="fim_image_upload">
    <el-dialog width="30%" title="图片上传" v-model="imgVisible" modal-class="image_dialog">
      <div class="image_preview_body">
        <el-image
            class="preview_img"
            :src="selectImg.src"
            :zoom-rate="1.2"
            :max-scale="7"
            :min-scale="0.2"
            :preview-src-list="[selectImg.src]"
            :initial-index="4"
            fit="cover"
        />
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
    <input type="file" ref="imgIptRef" @change="imageIptChange" style="display: none" accept="image/*">
    <slot><i title="图片上传" class="iconfont icon-tupian" @click="imgClick"></i></slot>
  </div>

</template>

<style lang="scss">
.image_dialog {
  .image_preview_body {
    display: flex;
    align-items: center;
    background-color: var(--item_hover);
    position: relative;
    padding: 20px;
    border-radius: 5px;

    .preview_img {
      width: 120px;
      height: 80px;
      object-fit: cover;
      border-radius: 5px;
      margin-right: 10px;
    }

    .info {

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