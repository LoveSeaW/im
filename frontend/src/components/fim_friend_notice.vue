<script setup lang="ts">
import {Edit} from "@element-plus/icons-vue";
import {ref} from "vue";
import {friendUpdateNicknameApi} from "@/api/user_api";
import {ElMessage} from "element-plus";

interface Props {
  friendId: number
  modelValue: string
}

const props = defineProps<Props>()
const emits = defineEmits(["update:modelValue"])
const isShow = ref(false)
const elIpt = ref()
const oldValueRef = ref("")

function iptFocus(oldValue: string) {
  isShow.value = true
  notice.value = oldValue
  oldValueRef.value = oldValue
  setTimeout(() => {
    elIpt.value?.focus()
  }, 100)
}

const notice = ref("")

async function iptBlur() {
  isShow.value = false
  if (notice.value !== oldValueRef.value) {
    let res = await friendUpdateNicknameApi(props.friendId, notice.value)
    if (res.code) {
      ElMessage.error(res.msg)
      return
    }
    ElMessage.success("更新好友备注成功")
    emits("update:modelValue", notice.value)
  }
}

</script>

<template>
  <div class="fim_friend_notice">

    <el-input
        v-model="notice"
        @blur="iptBlur"
        v-if="isShow"
        ref="elIpt"
        placeholder="修改好友备注"></el-input>
    <span v-else>{{ props.modelValue }}</span>
    <el-icon @click="iptFocus(props.modelValue)" :size="18">
      <Edit></Edit>
    </el-icon>
  </div>
</template>

<style lang="scss">
.fim_friend_notice {
  color: var(--text_color);
  display: flex;
  align-items: center;

  .el-icon {
    cursor: pointer;
    margin-left: 5px;
  }
}
</style>