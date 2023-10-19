<!--
 * @Author: SpenserCai
 * @Date: 2023-10-04 00:04:25
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-19 22:00:15
 * @Description: file content
-->
<script setup>
import { computed } from 'vue'
import { useMainStore } from '@/stores/main'
// import { mdiCheckDecagram } from '@mdi/js'
import BaseLevel from '@/components/BaseLevel.vue'
import { Button } from 'flowbite-vue'
import { mdiLock,mdiEarth } from '@mdi/js';
import UserAvatarCurrentUser from '@/components/UserAvatarCurrentUser.vue'
import CardBox from '@/components/CardBox.vue'
// import PillTag from '@/components/PillTag.vue'

const mainStore = useMainStore()

const userName = computed(() => mainStore.userName)
</script>

<template>
  <CardBox>
    <BaseLevel type="flex w-full items-center justify-between gap-4 justify-self-start">
      <div class="flex justify-start">
        <UserAvatarCurrentUser class="lg:mx-12" />
        <div class="space-y-3 text-center mt-2 md:text-left lg:mx-12">
          <h1 class="text-3xl">
            <b>{{ userName }}</b>
          </h1>
          <p class="text-slate-300">Created {{ mainStore.userImageCount }} Images</p>
          <p class="text-gray-500">
            <span class="text-gray-400">Joined on</span> {{ mainStore.created }}
          </p>        
        </div>
      </div>
      <div class="flex-1 mr-5">
        <Button v-if="mainStore.userIsPrivate" size="lg" gradient="purple-pink" class="text-lg font-bold uppercase" style="pointer-events:none">
          PRIVATE
          <template #prefix >
            <!--<BaseIcon :path="mdiLock" size="14" :hidden="true"/>-->
            <svg class="mr-1 -ml-2 w-5 h-5" fill="currentColor" viewBox="0 0 22 22" xmlns="http://www.w3.org/2000/svg">
              <path fill-rule="evenodd" :d="mdiLock" clip-rule="evenodd"></path>
            </svg>
          </template>
        </Button>
        <Button v-else size="lg" gradient="green-blue" class="text-lg font-bold uppercase" style="pointer-events:none">
          PUBLIC
          <template #prefix>
            <!--<BaseIcon :path="mdiEarth" size="14" :hidden="true"/>-->
            <svg class="mr-1 -ml-2 w-5 h-5" fill="currentColor" viewBox="0 0 22 22" xmlns="http://www.w3.org/2000/svg">
              <path fill-rule="evenodd" :d="mdiEarth" clip-rule="evenodd"></path>
            </svg>
          </template>
        </Button>
      </div>
    </BaseLevel>
  </CardBox>
</template>
