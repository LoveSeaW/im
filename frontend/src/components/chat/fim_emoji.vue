<script setup lang="ts">
import {ref} from "vue";
import {useStore} from "@/stores";

const store = useStore()
const menu = ref(1)

function sendUnicodeEmoji(e: Event) {
  const obj = e.target as HTMLElement
  if (obj.tagName === "SPAN") {
    store.setChatContent(obj.innerHTML)
  }
}

function sendIconfontEmoji(e: Event) {
  const obj = e.target as HTMLElement
  if (obj.tagName === "I") {
    obj.classList.add("iconfont_emoji")
    store.setChatContent(obj.outerHTML)
  }
}

function sendPackEmoji(e: Event) {
  const obj = e.target as HTMLElement
  if (obj.tagName === "IMG") {
    obj.classList.add("pack_emoji")
    store.setChatContent(obj.outerHTML)
  }
}

</script>

<template>
  <div class="fim_emoji">
    <el-popover popper-class="fim_emoji_popover"
                placement="top" trigger="click" :width="320">
      <template #reference>
        <i title="表情" class="iconfont icon-biaoqing"></i>
      </template>
      <template #default>
        <el-tabs tab-position="bottom" v-model="menu">
          <el-tab-pane class="classical" label="经典表情" :name="1">
            <div class="classical_emoji_list" @click="sendIconfontEmoji">
              <i class="iconfont icon-kaixin"></i>
              <i class="iconfont icon-face16"></i>
              <i class="iconfont icon-xieyanxiao"></i>
              <i class="iconfont icon-zhayanjing"></i>
              <i class="iconfont icon-bulini"></i>
              <i class="iconfont icon-wuyu"></i>
              <i class="iconfont icon-a-ziyuan101"></i>
            </div>
          </el-tab-pane>
          <el-tab-pane class="symbol" label="符号表情" :name="2">
            <div class="emoji_list" @click="sendUnicodeEmoji">
              <span title="嘿嘿">😀</span><span title="嘻嘻">😁</span><span title="笑哭了">😂</span>
              <span title="哈哈">😃</span><span title="大笑">😄</span><span title="苦笑">😅</span>
              <span title="斜眼笑">😆</span><span title="微笑天使">😇</span><span title="恶魔微笑">😈</span>
              <span title="眨眼">😉</span><span title="羞涩微笑">😊</span><span title="好吃">😋</span>
              <span title="松了口气">😌</span><span title="花痴">😍</span><span title="墨镜笑脸">😎</span>
              <span title="得意">😏</span><span title="冷漠">😐</span><span title="无语">😑</span>
              <span title="不高兴">😒</span><span title="汗">😓</span><span title="沉思">😔</span>
              <span title="困扰">😕</span><span title="困惑">😖</span><span title="亲亲">😗</span>
              <span title="飞吻">😘</span><span title="微笑亲亲">😙</span><span title="羞涩亲亲">😚</span>
              <span title="吐舌">😛</span><span title="单眼吐舌">😜</span><span title="眯眼吐舌">😝</span>
              <span title="失望">😞</span><span title="担心">😟</span><span title="生气">😠</span>
              <span title="怒火中烧">😡</span><span title="哭">😢</span><span title="痛苦">😣</span>
              <span title="傲慢">😤</span><span title="失望但如释重负">😥</span><span title="啊">😦</span>
              <span title="极度痛苦">😧</span><span title="害怕">😨</span><span title="累死了">😩</span>
              <span title="困">😪</span><span title="累">😫</span><span title="龇牙咧嘴">😬</span>
              <span title="放声大哭">😭</span><span title="吃惊">😮</span><span title="缄默">😯</span>
              <span title="冷汗">😰</span><span title="吓死了">😱</span><span title="震惊">😲</span>
              <span title="脸红">😳</span><span title="睡着了">😴</span><span title="晕头转向">😵</span>
              <span title="沉默">😶</span><span title="感冒">😷</span><span title="微微不满">🙁</span>
              <span title="呵呵">🙂</span><span title="倒脸">🙃</span><span title="翻白眼">🙄</span>
              <span title="微笑的猫">😸</span><span title="笑出眼泪的猫">😹</span><span title="大笑的猫">😺</span>
              <span title="花痴的猫">😻</span><span title="奸笑的猫">😼</span><span title="亲亲猫">😽</span>
              <span title="生气的猫">😾</span><span title="哭泣的猫">😿</span><span title="疲倦的猫">🙀</span>
              <span title="禁止手势">🙅</span><span title="OK手势">🙆</span><span title="鞠躬">🙇</span>
              <span title="非礼勿视">🙈</span><span title="非礼勿听">🙉</span><span title="非礼勿言">🙊</span>
              <span title="举手">🙋</span><span title="举双手">🙌</span><span title="皱眉">🙍</span>
              <span title="撅嘴">🙎</span><span title="双手合十">🙏</span>
            </div>
          </el-tab-pane>
          <el-tab-pane class="emoji_pack" label="表情包" :name="3">
            <div class="pack_emoji_list" @click="sendPackEmoji">
              <span><img src="@/assets/gif/QQ图片20240420004934.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420004938.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420004946.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420004950.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420004953.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420004956.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420005000.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420005004.gif" alt=""></span>
              <span><img src="@/assets/gif/QQ图片20240420005010.gif" alt=""></span>
            </div>
          </el-tab-pane>
        </el-tabs>
      </template>
    </el-popover>
  </div>
</template>

<style lang="scss">
.fim_emoji_popover {
  padding: 0;

  .symbol {
    .emoji_list {
      span {
        display: inline-flex;
        width: 25px;
        height: 25px;
        justify-content: center;
        align-items: center;
        cursor: pointer;

        &:hover {
          background-color: var(--item_hover);
        }
      }
    }
  }

  .classical {
    .classical_emoji_list {
      i {
        display: inline-flex;
        width: 40px;
        height: 40px;
        font-size: 25px;
        justify-content: center;
        align-items: center;
        cursor: pointer;

        &:hover {
          background-color: var(--item_hover);
        }
      }
    }
  }

  .emoji_pack {
    .pack_emoji_list {
      padding: 20px 10px 10px 10px;
      span {
        display: inline-flex;
        width: 70px;
        height: 70px;
        justify-content: center;
        align-items: center;
        cursor: pointer;
        transition: all .3s;

        &:hover {
          background-color: var(--item_hover);
          transform: scale(1.3);
          z-index: 1;
        }

        img {
          width: 68px;
          height: 68px;
          object-fit: cover;
        }
      }
    }
  }

  .el-tabs__header.is-bottom {
    margin-bottom: 0;
    margin-top: 0;
  }
}
</style>