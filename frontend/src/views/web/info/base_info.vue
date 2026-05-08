<script setup lang="ts">
import {useStore} from "@/stores";
import {ref, nextTick, watch} from "vue";
import {type userConfUpdateType, userInfoUpdateApi} from "@/api/user_api";
import {ElMessage} from "element-plus";

const store = useStore()

interface iptType {
  label: string
  val: boolean
  key: "friendOnline" | "sound"
  help: string
}

const list = ref<iptType[]>([
  {
    label: "提醒",
    val: store.userConfInfo.friendOnline,
    key: "friendOnline",
    help: "开启好友上线提醒"
  },
  {
    label: "声音",
    val: store.userConfInfo.sound,
    key: "sound",
    help: "关闭所有声音"
  }
])
const editRefList = ref()

watch(() => store.userConfInfo, () => {
  list.value[0].val = store.userConfInfo.friendOnline
  list.value[1].val = store.userConfInfo.sound
  if (store.userConfInfo.recallMessage) {
    isRecallMessage.value = true
    recallMessage.value = store.userConfInfo.recallMessage
  }
}, {deep: true})

async function change(index: number) {
  let data: userConfUpdateType = {}

  data[list.value[index].key] = list.value[index].val
  let res = await userInfoUpdateApi(data)
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success(list.value[index].label + "设置修改成功")
}


const isRecallMessage = ref(false)
const recallMessage = ref("撤回了一条消息")

async function recallMessageUpdate(val: boolean) {
  if (val){
    return
  }
  if (!isRecallMessage.value) {
    // 取消设置
    recallMessage.value = ""
  }
  let res = await userInfoUpdateApi({
    recallMessage: recallMessage.value,
  })
  if (res.code) {
    ElMessage.error(res.msg)
    return
  }
  ElMessage.success("消息会话设置成功")
}

</script>

<template>
  <div class="my_info_view">
    <el-form-item label="会话">
      <div>
        <div>
          <el-checkbox v-model="isRecallMessage" @change="recallMessageUpdate" label="撤回消息设置"></el-checkbox>
        </div>
        <div>
          <el-input v-if="isRecallMessage" style="width: 300px" @blur="recallMessageUpdate(false)" :maxlength="32"
                    placeholder="撤回消息提示词"
                    v-model="recallMessage"></el-input>
        </div>
      </div>


    </el-form-item>
    <el-form-item :label="item.label" v-for="(item, index) in list">
      <el-checkbox v-model="item.val" :label="item.help" @change="change(index)"></el-checkbox>
    </el-form-item>
  </div>
</template>