<script setup lang="ts">

interface Props {
  name: string
}
const props = defineProps<Props>()

const iconFileMap = {
  "pdf": "icon-pdf",
  "xlsx?": "icon-excel",
  "docx?": "icon-WORD",
  "psd": "icon-psd",
  "exe": "icon-exe",
  "xmind": "icon-Xmind",
  "zip|rar": "icon-zip-1",
  "jpe?g|png": "icon-img-cs",
  "txt": "icon-txtwenjian",
  "mp4": "icon-shipin_cs",
  "": "icon-wenjian-cs",
}

function getIconName(fileName: string): string {
  const fileList = fileName.split(".")
  const shuffix = fileList[fileList.length - 1]
  for (const shuffixKey in iconFileMap) {
    const regex = new RegExp(shuffixKey)
    if (regex.test(shuffix)) {
      return iconFileMap[shuffixKey as keyof typeof iconFileMap]
    }
  }
  return iconFileMap[""]

}
</script>

<template>
  <i class="iconfont" :class="getIconName(props.name)"></i>
</template>