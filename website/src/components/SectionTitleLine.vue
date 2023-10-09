<script setup>

import BaseIcon from '@/components/BaseIcon.vue'
import IconRounded from '@/components/IconRounded.vue'

defineProps({
  icon: {
    type: String,
    default: null
  },
  title: {
    type: String,
    required: true
  },
  main: Boolean
})

</script>

<template>
  <section :class="{ 'pt-6': !main }" class="mb-6 flex items-center justify-between">
    <div class="flex items-center justify-start">
      <IconRounded v-if="icon && main" :icon="icon" color="light" class="mr-3" bg />
      <BaseIcon v-else-if="icon" :path="icon" class="mr-2" size="20" />
      <h1 :class="main ? 'text-3xl' : 'text-xl'" class="leading-tight">
        {{ title }}
      </h1>
    </div>
  </section>
</template>
