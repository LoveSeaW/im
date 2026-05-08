<script setup lang="ts">
import {ChatLineRound, Menu, Notification, User} from "@element-plus/icons-vue";
import {useRouter} from "vue-router";
import {useStore} from "@/stores";

const store = useStore()
const router = useRouter()

function checkMenu(menu: string) {
  router.push({
    name: menu,
  })
}


</script>

<template>
  <div class="fim_slide">
    <div class="avatar" @click="checkMenu('info')">
      <img :src="store.userInfo.avatar" alt="">
    </div>
    <div class="fim_menus">
      <router-link class="icon" :to="{name: 'session_welcome'}">
        <el-icon>
          <ChatLineRound/>
        </el-icon>
      </router-link>
      <router-link class="icon" :to="{name: 'contacts'}">
        <el-icon>
          <User/>
        </el-icon>
      </router-link>
      <router-link class="icon" :to="{name: 'notice'}">
        <el-icon>
          <Notification/>
        </el-icon>
      </router-link>
    </div>

    <router-link class="other icon" :to="{name: 'info'}">
      <el-icon>
        <Menu/>
      </el-icon>
    </router-link>
  </div>
</template>

<style lang="scss">
.fim_slide {
  width: 80px;
  height: 100%;
  background-color: var(--slide_bg);
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 0;

  .avatar {
    width: 46px;
    height: 46px;

    img {
      width: 100%;
      height: 100%;
      border-radius: 50%;
    }
  }

  .fim_menus {
    margin-top: 20px;
    display: flex;
    flex-direction: column;
    align-items: center;

    > a {
      margin-bottom: 20px;
    }
  }

  .icon {
    font-size: 20px;
    width: 40px;
    height: 40px;
    border-radius: 5px;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: all .3s;

    &:hover {
      background-color: var(--icon_bg);
    }

    &.router-link-active {
      color: var(--active);
      background-color: var(--icon_bg);
    }
  }

  .other {
    position: absolute;
    bottom: 20px;
  }
}
</style>