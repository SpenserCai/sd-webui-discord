<script setup>
import { containerMaxW } from '@/config.js'
import BaseLevel from '@/components/BaseLevel.vue'

</script>

<template>
  <footer class="py-2 px-6" :class="containerMaxW">
    <BaseLevel>
      <div class="text-center md:text-left">
        <b>Powered by <a href="https://github.com/SpenserCai/sd-webui-discord" target="_blank">SD-WEBUI-DISCORD</a></b>
        <slot />
      </div>
      
    </BaseLevel>
  </footer>
</template>
