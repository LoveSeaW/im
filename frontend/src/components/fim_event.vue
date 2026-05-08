<script setup lang="ts">
import {watch} from "vue";
import type {eventKeyType} from "@/types/type";
import {eventStore} from "@/stores/event_store";

const store = eventStore()

interface Props {
  eventKey: eventKeyType
}

const emits = defineEmits(["event"])
const props = defineProps<Props>()

watch(() => store.eventKey, () => {
  if (store.eventKey === props.eventKey) {
    emits("event")
  }
})

</script>

<template>

</template>

<style scoped lang="scss">

</style>