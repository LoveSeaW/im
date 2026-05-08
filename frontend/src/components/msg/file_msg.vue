<script setup lang="ts">
import {fileSizeFormat} from "@/utils/file";
import Fim_file_icon from "@/components/chat/fim_file_icon.vue";
import type {fileMsg} from "@/types/msg";

interface Props {
  data: fileMsg
  isMe?: boolean
}

const props = defineProps<Props>()

</script>

<template>
  <div class="file_msg">
    <div class="file_top">
      <fim_file_icon :name="props.data.title as string"></fim_file_icon>
      <div class="info">
        <div class="name">{{ props.data.title }}</div>
        <div class="size">{{ fileSizeFormat(props.data.size as number) }}</div>
      </div>
    </div>
    <div class="file_action">

      <a :href="props.data.src" :download="props.data.title">
        <el-button size="small" type="text">
          下载
        </el-button>
      </a>
    </div>
  </div>
</template>

<style lang="scss">
.file_msg {
  border-radius: 5px;
  border: 1px solid var(--border_color);
  min-width: 200px;

  .file_top {
    padding: 20px 10px;
    border-bottom: 1px solid var(--border_color);
    display: flex;
    align-items: center;

    .iconfont {
      font-size: 30px;
    }

    .info {
      margin-left: 10px;
    }

    .name {

    }

    .size {
      font-size: 12px;
      margin-top: 5px;
      color: var(--text_color);
    }
  }

  .file_action {
    padding: 10px;
    display: flex;
    justify-content: end;
  }
}
</style>