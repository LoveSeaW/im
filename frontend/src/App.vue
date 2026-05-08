<script setup lang="ts">
import {useStore} from "@/stores";

const store = useStore()
store.loadToken()
store.loadTheme()
store.getSettingsInfo()
store.initChatWs()
store.initGroupWs()
</script>

<template>
  <div class="fim_body">
    <router-view/>
  </div>
</template>

<style>
.fim_body {
  width: 100%;
  height: 100vh;
  overflow-y: hidden;
  background-color: var(--bg);
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
