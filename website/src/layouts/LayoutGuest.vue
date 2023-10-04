<!--
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-04 11:54:51
 * @Description: file content
-->
<script setup>
import { useDarkModeStore } from '@/stores/darkMode.js'
const darkModeStore = useDarkModeStore()
darkModeStore.set(true)
</script>
<template>
  <div class="bg-gray-50 dark:bg-slate-800 dark:text-slate-100">
    <slot />
  </div>
</template>
