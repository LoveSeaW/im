<script setup lang="ts">
import {ref, watch} from "vue";
import {uploadImageApi} from "@/api/file_api";
import {ElMessage} from "element-plus";
import {base64UrlToFile} from "@/utils/file";
import {useStore} from "@/stores";
import {dateTimeFormat} from "@/utils/date";
import type {msgChatType} from "@/types/msg";
import {Msg} from "@/service/send_msg";

const store = useStore()

interface Props {
  id: number
  type: "user" | "group"
}

const props = defineProps<Props>()
const msg = new Msg(props.type)
const paste = async (event: ClipboardEvent) => {

  const items = (event.clipboardData as DataTransfer).items;
  for (const index in items) {
    const item = items[index];
    if (item.kind === 'file') {
      // 取消默认事件
      event.preventDefault();
      const blob = item.getAsFile();
      // 上传图片

      // const res = await uploadImageApi(blob as File, "chat")
      // if (res.code){
      //   ElMessage.error(res.msg)
      //   continue
      // }

      // insertImageByLink(res.data.url)
      const reader = new FileReader();
      reader.onload = (event: ProgressEvent<FileReader>) => {
        // 这里是在前端显示图片，可以自由操作是上传后端还是进行文件操作
        insertImageByLink((event.target as FileReader).result as string)
        return
      };
      reader.readAsDataURL(blob as File);
    }
  }
}

function insertImageByLink(link: string) {
  const image = new Image()
  image.src = link
  const selection = window.getSelection() as Selection
  const range = selection.getRangeAt(0)
  range.insertNode(image)

}


const iptRef = ref<HTMLDivElement>()

async function uploadImage(file: File) {
  const res = await uploadImageApi(file, "chat")
  return `<img class="paste_img" src="${res.data.url}">`
}

async function send() {


  // 前端先根据文本内容中的img  i 这两个标签，判断是不是图文消息
  if (/<img .*?>|<i.*?>[\s\S]*?<\/i>/.test(iptRef.value?.innerHTML as string)) {
    // 图文消息
    const promiseList: any[] = []
    const regex = /<img src="(.*?)">/g
    iptRef.value?.innerHTML.replace(regex, function (substring: string, ...args: string[]): string {
      if (args[0].startsWith("data:image/png;base64")) {
        const file = base64UrlToFile(args[0], "img.png")
        promiseList.push(uploadImage(file))
        return ""
      }

      return ""
    })

    let newValue = iptRef.value?.innerHTML as string
    if (promiseList.length) {
      // 有图片
      const data = await Promise.all(promiseList)
      newValue = iptRef.value?.innerHTML.replace(regex, (substring: string, ...args: string[]): string => {
        if (args[0].startsWith("data:image/png;base64")) {
          return data.shift()
        }
        return substring
      }) as string
    }
    msg.imageTextMsg(props.id as number, newValue);
    (document.querySelector(".editor") as HTMLDivElement).innerHTML = ""

    return
  }

  if (replyMsgVisible.value) {
    msg.replyMsg(props.id as number, {
      msgID: store.replayChatMsgData?.id as number,
      content: iptRef.value?.innerText as string,
    })
    store.replayChatMsgData = undefined
  } else {
    msg.textMsg(props.id as number, iptRef.value?.innerText as string);
  }
  (document.querySelector(".editor") as HTMLDivElement).innerHTML = ""


}

watch(() => store.chatContent, () => {
  (document.querySelector(".editor") as HTMLDivElement).innerHTML += store.chatContent
})

watch(() => store.replayChatMsgData, () => {
  if (store.replayChatMsgData?.id) {
    replyMsgVisible.value = true
    iptRef.value?.focus()
  } else {
    replyMsgVisible.value = false
  }
}, {deep: true})


const replyMsgVisible = ref(false)

function replyRemove() {
  store.replayChatMsgData = undefined
}

function keyDown(e: KeyboardEvent) {
  if (e.key === "Backspace") {
    // 按下退格按键
    // 当内容为空，并且有显示回复内容的时候。按下退格就是删除回复内容
    if (iptRef.value?.innerText === "" && replyMsgVisible.value) {
      store.replayChatMsgData = undefined
    }
  }
}


</script>

<template>
  <div class="fim_chat_ipt">
    <div class="reply_msg" v-if="replyMsgVisible">
      <div class="reply_inner_msg">
        <div class="head">
          <span class="nickname">{{ store.replayChatMsgData?.user.nickName }}</span>
          <span class="date">{{ dateTimeFormat((store.replayChatMsgData as msgChatType).created_at) }}</span>
        </div>
        <div class="body">
          {{ store.replayChatMsgData?.msgPreview }}
        </div>

        <div class="close" @click="replyRemove">
          <i class="iconfont icon-quxiao"></i>
        </div>

      </div>
    </div>
    <div class="editor" @keydown="keyDown" ref="iptRef" @paste="paste" contenteditable="true"></div>
    <div class="btn">
      <el-button @click="send" size="small" type="primary">发送</el-button>
    </div>
  </div>
</template>

<style lang="scss">
.fim_chat_ipt {
  position: relative;

  .btn {
    position: absolute;
    right: 10px;
    bottom: 10px;
  }

  .reply_msg {
    padding: 0 5px;

    .reply_inner_msg {
      background-color: var(--item_hover);
      padding: 5px;
      border-radius: 5px;
      font-size: 12px;
      color: var(--text_color);
      position: relative;

      .head {
        margin-bottom: 5px;
      }

      .nickname {
        margin-right: 5px;
      }

      .close {
        position: absolute;
        right: 5px;
        width: 20px;
        height: 20px;
        border-radius: 50%;
        background-color: var(--border_color);
        display: flex;
        justify-content: center;
        align-items: center;
        top: 50%;
        transform: translateY(-50%);
        cursor: pointer;

        i {
          font-size: 12px;
        }
      }
    }

    & + .editor {
      height: 72px;
    }
  }

  .editor {
    height: 120px;
    border-radius: 5px;
    padding: 5px;
    overflow-y: auto;

    &:focus-visible {
      outline: none;
    }

    /* 整个滚动条 */
    &::-webkit-scrollbar {
      width: 12px; /* 滚动条的宽度 */
      height: 12px; /* 滚动条的高度，对水平滚动条有效 */
      background-color: #f9f9fd; /* 滚动条的背景颜色 */
    }

    /* 滚动条轨道 */
    &::-webkit-scrollbar-track {
      border-radius: 10px;
      background: #e1e1e1; /* 轨道的背景颜色 */
    }

    /* 滚动条滑块 */
    &::-webkit-scrollbar-thumb {
      border-radius: 10px;
      background-color: #c1c1c1; /* 滑块的背景颜色 */
      border: 3px solid #e1e1e1; /* 滑块的边框和轨道相同的颜色，可以制造“边距”的效果 */
    }

    /* 滚动条滑块：悬停效果 */
    &::-webkit-scrollbar-thumb:hover {
      background-color: #a8a8a8; /* 滑块的悬停颜色 */
    }

    /* 滚动条滑块：激活时的效果 */
    &::-webkit-scrollbar-thumb:active {
      background-color: #888888; /* 滑块的激活颜色 */
    }

    /* 滚动条按钮（上下箭头） */
    &::-webkit-scrollbar-button {
      display: none; /* 通常情况下不显示滚动条按钮 */
    }

  }
}
</style>