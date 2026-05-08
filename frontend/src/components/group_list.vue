<script setup lang="ts">
import type {groupType} from "@/api/group_api";
import router from "@/router";
import {useRoute} from "vue-router";

const route = useRoute()
interface Props {
  list: groupType[]
}

const props = defineProps<Props>()

function goGroup(record: groupType) {
  router.push({
    name: "group_detail",
    params: {
      id: record.groupId,
    }
  })
}

</script>

<template>
  <div class="group_list_com">
    <div class="item" :class="{active: Number(route.params.id) === item.groupId && route.name === 'group_detail'}" @click="goGroup(item)" v-for="item in props.list">
      <div class="avatar">
        <img v-if="item.groupAvatar.startsWith('/api')" :src="item.groupAvatar" alt="">
        <span v-else>{{ item.groupAvatar }}</span>
      </div>
      <div class="info">
        <div class="title">
          <el-text style="max-width: 11rem" truncated>{{ item.groupTitle }}</el-text>
        </div>
        <div class="people">{{ item.groupMemberCount }} 人</div>
      </div>
    </div>
  </div>
</template>

<style lang="scss">
.group_list_com {
  width: 100%;

  .item {
    height: 50px;
    display: flex;
    padding: 0 10px;
    align-items: center;
    cursor: pointer;

    &:hover {
      background-color: var(--item_hover);
    }

    &.active {
      background-color: var(--item_hover);
    }


    .avatar {
      width: 40px;
      display: flex;
      align-items: center;

      img {
        width: 35px;
        height: 35px;
        border-radius: 5px;
        object-fit: cover;
      }

      span {
        width: 35px;
        height: 35px;
        border-radius: 5px;
        color: white;
        background-color: #529b2e;
        display: flex;
        justify-content: center;
        align-items: center;
        overflow: hidden;
      }
    }

    .info {
      font-size: 14px;

      .title {
        font-weight: 600;
      }

      .people {
        color: #555;
      }
    }
  }
}
</style>