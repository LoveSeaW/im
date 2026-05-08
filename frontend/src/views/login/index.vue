<script setup lang="ts">
import Qq_color from "@/components/fim_login/qq_color.vue";
import {useStore} from "@/stores";
import {authOpenLoginApi} from "@/api/auth_api";
import {ElMessage} from "element-plus";
import router from "@/router";
import {useRoute} from "vue-router";
const route = useRoute()
const store = useStore()

function gotoQQLogin() {
  window.open(store.settingsInfo.qq?.webPath, "_self")
}

async function initRouter(){
  const flag = route.query.flag
  const code = route.query.code
  if (flag && code){
    let res = await authOpenLoginApi({
      code: code as string,
      flag: flag as string
    })
    if (res.code){
      ElMessage.error(res.msg)
      return
    }
    ElMessage.success("登陆成功")
    store.setToken(res.data.token)
    router.push({
      name: "web",
    })
  }
}
initRouter()

</script>

<template>
  <div class="fim_login">
    <div class="banner">
      <qq_color></qq_color>
    </div>
    <div class="login_form">
      <router-view></router-view>
      <div class="other_login">
        <div class="label">第三方登陆</div>
        <div class="icons">
          <i class="iconfont icon-QQ" v-if="store.settingsInfo.qq?.enable" @click="gotoQQLogin"></i>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.fim_login {
  width: 500px;
  background-color: white;
  border-radius: 5px;
  overflow: hidden;
  box-shadow: 0 0 5px 3px rgba(0, 0, 0, 0.05);

  .banner {
    height: 140px;
    width: 100%;
    background-color: #d9d9d9;
  }

  .login_form {
    padding: 20px 80px;

    .item_password, .item_action, .item_btn {
      margin-bottom: 6px;
    }

    .other_login {
      display: flex;
      flex-direction: column;
      align-items: center;

      .label {
        font-size: 14px;
        color: #555;
        display: flex;
        align-items: center;
        width: 100%;
        justify-content: space-between;

        &::before, &::after {
          width: 35%;
          height: 1px;
          background-color: #e3e3e3;
          content: "";
          display: inline-flex;
        }
      }

      .icons {
        margin-top: 5px;

        i {
          font-size: 36px;
          cursor: pointer;
        }
      }
    }
  }

}
</style>