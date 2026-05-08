<script setup lang="ts">
import type {imageTextMsg} from "@/types/msg";
import {createApp, h, onMounted} from "vue";
import {ElImage} from "element-plus";

interface Props {
  data: imageTextMsg
}

const props = defineProps<Props>()

onMounted(() => {
  const imgList = document.querySelectorAll(".image_text_msg .paste_img")
  imgList.forEach((value: Node) => {
    const src = (value as HTMLImageElement).src
    const elImageVNode = h(ElImage, {
      src: src,
      zoomRate: 1.2,
      maxScale: 7,
      minScale: 0.2,
      previewSrcList: [src],
      initialIndex: 4,
      fit: "cover"
    })
    const app = createApp(() => elImageVNode)
    app.mount(value.parentNode as Element)
  })
})

</script>

<template>
  <div class="image_text_msg">
    <div v-html="props.data.content"></div>
  </div>

</template>

<style lang="scss">
.image_text_msg {
  border-radius: 5px;
  padding: 5px;
  color: var(--text_color);
  background-color: var(--item_hover);

  img {
    max-width: 300px;
  }

  .pack_emoji {
    width: 60px;
    height: 60px;
    object-fit: cover;
  }

  .iconfont_emoji {
    font-size: 20px;
  }
}
</style>