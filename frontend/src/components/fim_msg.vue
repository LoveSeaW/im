<script setup lang="ts">
import Fim_avatar from "@/components/fim_avatar.vue";
import {createApp, h, onMounted, ref, watch} from "vue";
import type {
  fileMsg,
  imageMsg,
  imageTextMsg,
  msgChatType,
  replyMsg,
  videoCallMsg,
  videoMsg,
  voiceMsg, withdrawMsg
} from "@/types/msg";
import Image_text_msg from "@/components/msg/image_text_msg.vue";
import Video_msg from "@/components/msg/video_msg.vue";
import Voice_msg from "@/components/msg/voice_msg.vue";
import File_msg from "@/components/msg/file_msg.vue";
import Image_msg from "@/components/msg/image_msg.vue";
import Video_call_msg from "@/components/msg/video_call_msg.vue";
import {copyText} from "@/utils/copy";
import {useStore} from "@/stores";
import Reply_msg from "@/components/msg/reply_msg.vue";
import {dateTimeFormat, timeFormat} from "../utils/date";
import {ElMessage} from "element-plus";
import {eventStore} from "@/stores/event_store";
import Withdraw_msg from "@/components/msg/withdraw_msg.vue";
import {useRoute} from "vue-router";
import type {baseResponse} from "@/api";
import {chatDeleteApi} from "@/api/chat_api";
import {groupChatDeleteApi} from "@/api/group_api";
import {Msg} from "@/service/send_msg";

const route = useRoute()
const eStore = eventStore()
const store = useStore()

interface Props {
  data: msgChatType
  showNickname?: boolean
  type: "user" | "group"
  checkShow?: boolean
}

const props = defineProps<Props>()
const emits = defineEmits(["checkShow", "check"])
const msg = new Msg(props.type)

const {showNickname = false} = props

const isCheck = ref(false)

const checkShow = ref(false)

function showCheck(id: number) {
  // 通知父组件，让子组件每一个都显示复选框
  emits("checkShow")
  isCheck.value = true
}

watch(isCheck, () => {
  if (isCheck.value) {
    // 就向父组件 添加 当前的消息id
    emits("check", "add", props.data.id)
  } else {
    // 就向父组件 移除 当前的消息id
    emits("check", "remove", props.data.id)
  }
})

function closeCheck() {
  isCheck.value = false
}

function copy() {
  const dom = document.querySelector(`[msg_id="${props.data.id}"]`) as HTMLDivElement
  if (dom) {
    copyText(dom.innerText)
  }
}

function reply() {
  store.replayChatMsgData = props.data
}

function withdraw() {
  // 算一下消息的时间，2分钟以内才能撤回
  if (((new Date().getTime() - new Date(props.data.created_at).getTime()) / 1000) > 120) {
    ElMessage.warning("只能撤回两分钟以内的消息")
    return
  }

  msg.withdrawMsg(Number(route.params.id), props.data.id)
}

async function msgRemove() {
  let res: baseResponse<string> = {code: 0, msg: "", data: ""}
  if (props.type === "user") {
    res = await chatDeleteApi([props.data.id])
  } else {
    res = await groupChatDeleteApi(Number(route.params.id), [props.data.id])
  }
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }

  // 通知父组件刷新
  eStore.setEventKey("history")
}

defineExpose({
  closeCheck
})


</script>

<template>
  <div class="fim_msg" :class="[props.data.isMe?'isMe': 'isYou', props.type]">
    <div class="msg_date">
      <span class="middle_date" v-if="props.data.showDate">
          {{ props.data.created_at }}
      </span>
      <span class="msg_middle_date" v-if="store.showMsgDate">
          {{ timeFormat(props.data.created_at) }}
      </span>
    </div>
    <withdraw_msg v-if="props.data.msg.type === 8" :is-me="props.data.isMe"
                  :data="props.data"></withdraw_msg>
    <div class="msg_inner" v-else>
      <el-checkbox v-if="props.checkShow" v-model="isCheck"></el-checkbox>
      <fim_avatar :src="props.data.user.avatar"></fim_avatar>
      <div class="info">
        <el-text v-if="showNickname && !props.data.isMe" truncated style="max-width: 10rem" class="nickname">{{
            props.data.user.nickName
          }}
        </el-text>
        <el-dropdown trigger="contextmenu">
          <div class="msg" :class="'t_' + props.data.msg.type" :msg_id="props.data.id">
            <div class="text_msg" v-if="props.data.msg.type === 1">
              {{ props.data.msg.textMsg?.content }}
            </div>
            <image_msg v-else-if="props.data.msg.type === 2" :data="props.data.msg.imageMsg as imageMsg"></image_msg>
            <video_msg v-else-if="props.data.msg.type === 3" :data="props.data.msg.videoMsg as videoMsg"></video_msg>
            <file_msg v-else-if="props.data.msg.type === 4" :data="props.data.msg.fileMsg as fileMsg"></file_msg>
            <voice_msg :is-me="props.data.isMe" v-else-if="props.data.msg.type === 5"
                       :data="props.data.msg.voiceMsg as voiceMsg"></voice_msg>
            <video_call_msg v-else-if="props.data.msg.type === 7" :is-me="props.data.isMe"
                            :data="props.data.msg.videoCallMsg as videoCallMsg"></video_call_msg>
            <!--            <withdraw_msg v-else-if="props.data.msg.type === 8" :is-me="props.data.isMe"-->
            <!--                          :data="props.data.msg.withdrawMsg as withdrawMsg"></withdraw_msg>-->
            <reply_msg v-else-if="props.data.msg.type === 9" :data="props.data.msg.replyMsg as replyMsg"
                       :is-me="props.data.isMe"></reply_msg>
            <image_text_msg v-else-if="props.data.msg.type === 14"
                            :data="props.data.msg?.imageTextMsg as imageTextMsg"></image_text_msg>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="copy">复制</el-dropdown-item>
              <el-dropdown-item @click="reply">回复</el-dropdown-item>
              <el-dropdown-item
                  v-if="props.data.user.id === store.userInfo.userID &&  (((new Date().getTime() - new Date(props.data.created_at).getTime()) / 1000) < 120) "
                  @click="withdraw">撤回消息
              </el-dropdown-item>
              <el-dropdown-item @click="showCheck(props.data.id)">选择</el-dropdown-item>
              <el-dropdown-item style="color: red" @click="msgRemove">删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

      </div>
    </div>
  </div>
</template>

<style lang="scss">
.fim_msg {
  display: flex;
  width: 100%;
  flex-direction: column;

  .msg_date {
    width: 100%;
    font-size: 12px;
    color: var(--text_color);
    margin-bottom: 6px;
    margin-top: 2px;
    text-align: center;
    position: relative;
    height: 16px;

    .middle_date {
      text-align: center;
    }

    .msg_middle_date {
      position: absolute;
      right: 10px;
    }
  }

  .msg_inner {
    width: 90%;
    display: flex;

    .el-checkbox {
      margin-right: 5px;
    }

  }

  .info {
    margin-left: 10px;
    position: relative;

    .el-dropdown {
      display: flex;
    }


    .msg {
      display: flex;
      justify-content: start;


      &.t_1, &.t_5, &.t_7, &.t_14, &.t_9 {
        &::before {
          position: absolute;
          content: "";
          display: block;
          width: 0;
          height: 0;
          border: 5px solid transparent;
          border-right-color: var(--item_hover);
          top: 14px;
          left: -9px;
        }
      }


      .text_msg {
        font-size: 14px;
        color: var(--text_color);
        background-color: var(--item_hover);
        padding: 10px;
        border-radius: 5px;
        width: fit-content;
        min-height: 40px;
        white-space: break-spaces;
      }

      .preview_img {
        width: 120px;
        height: 80px;
        border-radius: 5px;
      }


    }
  }

  &.isMe {
    //justify-content: right;
    align-items: end;

    .msg_inner {
      flex-direction: row-reverse;

      .el-checkbox {
        margin-left: 5px;
      }
    }

    .info {
      margin-left: 0;
      margin-right: 10px;

      .msg {
        justify-content: end;

        &::before {
          left: inherit;
          right: -8px;
          border-right-color: transparent;
          border-left-color: var(--msg_active_color);
        }

        .text_msg {
          background-color: var(--msg_active_color);
          color: white;
        }

        .image_text_msg {
          background-color: var(--msg_active_color);
          color: white;
        }
      }
    }
  }

  &.group {
    &.isYou {
      .info {
        margin-top: -10px;

        &::before {
          top: 30px;
        }
      }
    }

  }
}
</style>