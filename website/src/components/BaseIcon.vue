<!--
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-07 17:33:36
 * @Description: file content
-->
<script setup>
import { computed } from 'vue'

const props = defineProps({
  path: {
    type: String,
    required: true
  },
  w: {
    type: String,
    default: 'w-6'
  },
  h: {
    type: String,
    default: 'h-6'
  },
  size: {
    type: [String, Number],
    default: null
  },
  class: {
    type: String,
    default: null
  }
})

const spanClass = computed(() => `inline-flex justify-center items-center ${props.w} ${props.h}`)

const iconSize = computed(() => props.size ?? 16)
</script>

<template>
  <span :class="spanClass">
    <svg viewBox="0 0 24 24" :width="iconSize" :height="iconSize" class="inline-block">
      <path fill="currentColor" :class="props.class" :d="path" />
    </svg>
    <slot />
  </span>
</template>
