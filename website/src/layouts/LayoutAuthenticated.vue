<!--
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-11 11:55:03
 * @Description: file content
-->
<script setup>
import { mdiForwardburger, mdiBackburger, mdiMenu } from '@mdi/js'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import menuAside from '@/menuAside.js'
import menuNavBar from '@/menuNavBar.js'
import { useDarkModeStore } from '@/stores/darkMode.js'
import BaseIcon from '@/components/BaseIcon.vue'
import FormControl from '@/components/FormControl.vue'
import NavBar from '@/components/NavBar.vue'
import NavBarItemPlain from '@/components/NavBarItemPlain.vue'
import AsideMenu from '@/components/AsideMenu.vue'
import NotifyGroup from '@/components/NotifyGroup.vue'
// import FooterBar from '@/components/FooterBar.vue'
import Cookies from "js-cookie"

const layoutAsidePadding = 'xl:pl-60'

const darkModeStore = useDarkModeStore()
darkModeStore.set(true)

const router = useRouter()

const isAsideMobileExpanded = ref(false)
const isAsideLgActive = ref(false)

router.beforeEach(() => {
  isAsideMobileExpanded.value = false
  isAsideLgActive.value = false
})

const menuClick = (event, item) => {
  if (item.isToggleLightDark) {
    darkModeStore.set()
  }

  if (item.isLogout) {
    Cookies.remove('token')
    router.push('/')
    
  }
}
</script>

<template>
  <div
    :class="{
      'overflow-hidden lg:overflow-visible': isAsideMobileExpanded
    }"
  >
    <div
      :class="[layoutAsidePadding, { 'ml-60 lg:ml-0': isAsideMobileExpanded }]"
      class="pt-14 min-h-screen w-screen transition-position lg:w-auto bg-gray-50 dark:bg-slate-800 dark:text-slate-100"
    >
      <NavBar
        :menu="menuNavBar"
        :class="[layoutAsidePadding, { 'ml-60 lg:ml-0': isAsideMobileExpanded }]"
        @menu-click="menuClick"
      >
        <NavBarItemPlain
          display="flex lg:hidden"
          @click.prevent="isAsideMobileExpanded = !isAsideMobileExpanded"
        >
          <BaseIcon :path="isAsideMobileExpanded ? mdiBackburger : mdiForwardburger" size="24" />
        </NavBarItemPlain>
        <NavBarItemPlain display="hidden lg:flex xl:hidden" @click.prevent="isAsideLgActive = true">
          <BaseIcon :path="mdiMenu" size="24" />
        </NavBarItemPlain>
        <NavBarItemPlain use-margin>
          <FormControl placeholder="Search (ctrl+k)" ctrl-k-focus transparent borderless />
        </NavBarItemPlain>
      </NavBar>
      <AsideMenu
        :is-aside-mobile-expanded="isAsideMobileExpanded"
        :is-aside-lg-active="isAsideLgActive"
        :menu="menuAside"
        @menu-click="menuClick"
        @aside-lg-close-click="isAsideLgActive = false"
      />

      <slot />
      <NotifyGroup group="authenticated" />
    </div>
  </div>
</template>
