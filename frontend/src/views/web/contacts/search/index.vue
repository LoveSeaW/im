<script setup lang="ts">
import {useRoute} from "vue-router";
import {useStore} from "@/stores";
import {ref} from "vue";

const store = useStore()
const route = useRoute()


const searchKey = ref("")

function search() {
  store.searchData = {
    name: route.name as "search_user" | "search_group" | "",
    value: searchKey.value,
  }
}

</script>

<template>
  <div class="fim_search">
    <div class="search_head">
      <el-input @keydown.enter="search" v-model="searchKey"
                :placeholder="route.name === 'search_user' ? '请输入用户号、用户昵称' : '请输入群号、群名称'"></el-input>
      <el-button type="primary" @click="search">搜索</el-button>
    </div>
    <div class="search_menu">
      <router-link :to="{name:'search_user'}">找人</router-link>
      <router-link :to="{name:'search_group'}">找群</router-link>
    </div>
    <div class="search_main">
      <router-view/>
    </div>
  </div>
</template>

<style lang="scss">
.fim_search {
  .search_head {
    padding: 40px 20px;
    border-bottom: 1px solid var(--border_color);
    display: flex;
    justify-content: center;

    .el-input {
      width: 300px;
      margin-right: 5px;
    }
  }

  .search_menu {
    height: 40px;
    border-bottom: 1px solid var(--border_color);
    display: flex;
    align-items: center;
    padding: 0 20px;

    a {
      margin-right: 20px;

      &.router-link-active {
        color: var(--active);
      }
    }
  }

  .search_main {
    padding: 20px;
  }
}
</style>