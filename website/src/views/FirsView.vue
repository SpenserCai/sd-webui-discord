<!--
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-04 10:41:04
 * @Description: file content
-->
<script setup>
import { gradientBgDark } from '@/colors.js'
import SectionMain from '@/components/SectionMain.vue'
import BaseButton from '@/components/BaseButton.vue'
import LayoutGuest from '@/layouts/LayoutGuest.vue'
import { mdiLogin,mdiGithub } from '@mdi/js'
import router from '@/router'

const Login = () => {
  router.push('/app')
}

const Github = () => {
  window.location.href = 'https://github.com/SpenserCai/sd-webui-discord'
}
</script>

<template>
  <LayoutGuest>
    <div :class="gradientBgDark" class="flex min-h-screen items-center justify-center">
      <SectionMain>
        <h1 class="text-4xl md:text-5xl text-center text-white font-bold mt-12 mb-3 lg:mt-0">
          SD-WEBUI-DISCORD
        </h1>
        <h2 class="text-xl md:text-xl text-center text-white mb-12">
          Support For Clustered Stable Diffusion WebUi Discord Bot
        </h2>
        <div class="grid gap-6 grid-cols-1 lg:grid-cols-2 px-6 max-w-6xl mx-auto">
          <!--一个Login按钮-->
          <BaseButton label="Login" color="info" class="w-full" :icon="mdiLogin" @click="Login"/>
          <BaseButton label="Github" color="success" class="w-full" :icon="mdiGithub" @click="Github"/>
        </div>
      </SectionMain>
    </div>
  </LayoutGuest>
</template>
