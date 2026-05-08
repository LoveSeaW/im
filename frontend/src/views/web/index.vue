<script setup lang="ts">
import Fim_slide from "@/components/fim_slide.vue";
import Fim_video_call_rev from "@/components/chat/fim_video_call_rev.vue";
</script>

<template>
  <div class="fim_web">
    <fim_slide></fim_slide>
    <fim_video_call_rev></fim_video_call_rev>
    <div class="fim_main">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss">
.fim_web {
  width: 1000px;
  height: 600px;
  background-color: var(--el-bg-color);
  border-radius: 5px;
  overflow: hidden;
  box-shadow: 0 0 5px 3px rgba(0, 0, 0, 0.05);
  display: flex;
}

.fim_main {
  width: calc(100% - 80px);
}
</style>