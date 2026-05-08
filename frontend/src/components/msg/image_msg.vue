<script setup lang="ts">
import {ElImage} from "element-plus";
import type {imageMsg} from "@/types/msg";
interface Props {
  data: imageMsg
  isMe?: boolean
}

const props = defineProps<Props>()

</script>

<template>
  <div class="image_msg">
    <el-image
        :title="props.data.title"
        class="preview_img"
        :src="props.data.src as string"
        :zoom-rate="1.2"
        :max-scale="7"
        :min-scale="0.2"
        :preview-src-list="[props.data.src as string]"
        :initial-index="4"
        fit="cover"
    />
  </div>
</template>

<style lang="scss">

</style>