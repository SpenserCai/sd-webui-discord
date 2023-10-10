<script setup>
import { computed } from 'vue'
import { useMainStore } from '@/stores/main'
// import { mdiCheckDecagram } from '@mdi/js'
import BaseLevel from '@/components/BaseLevel.vue'
import UserAvatarCurrentUser from '@/components/UserAvatarCurrentUser.vue'
import CardBox from '@/components/CardBox.vue'
// import PillTag from '@/components/PillTag.vue'

const mainStore = useMainStore()

const userName = computed(() => mainStore.userName)
</script>

<template>
  <CardBox>
    <BaseLevel type="flex lg:justify-start">
      <UserAvatarCurrentUser class="lg:mx-12" />
      <div class="space-y-3 text-center md:text-left lg:mx-12">
        
        <h1 class="text-3xl">
          <b>{{ userName }}</b>
        </h1>
        <p class="text-slate-300">Created {{ mainStore.userImageCount }} Images</p>
        <p class="text-gray-500">
          <span class="text-gray-400">Joined on</span> {{ mainStore.created }}
        </p>        
      </div>
    </BaseLevel>
  </CardBox>
</template>
