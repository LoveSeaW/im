<script setup lang="ts">

</script>

<template>
  <div class="fim_slide_head">
    <router-link :to="{name: 'search_user'}"><i class="iconfont icon-sousuo"></i>搜索</router-link>
    <router-link :to="{name: 'add_group'}"><i class="iconfont icon-chuangjianqunliao"></i>添加群聊</router-link>
  </div>
</template>

<style lang="scss">
.fim_slide_head {
  height: 40px;
  padding: 0 10px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border_color);
  color: var(--text_color);

  span {
    cursor: pointer;
  }

  i {
    margin-right: 5px;
    color: var(--text_color);
  }
}
</style>