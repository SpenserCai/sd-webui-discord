<!--
 * @Author: SpenserCai
 * @Date: 2023-10-03 16:06:09
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-04 01:09:55
 * @Description: file content
-->
<script setup>
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import UserCard from '@/components/UserCard.vue'
import SectionMain from '@/components/SectionMain.vue'
import CardBox from '@/components/CardBox.vue'
import { useMainStore } from '@/stores/main'
import { Input } from 'flowbite-vue'

const mainStore = useMainStore()
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <UserCard class="mb-6" />
        <CardBox is-form>
          <div class="grid grid-cols-2 gap-2">
            <Input
              v-model="mainStore.stableConfig.sd_model_checkpoint"
              label="Checkpoints"
              disabled>
            </Input>
            <Input
              v-model="mainStore.stableConfig.sd_vae"
              label="Vae"
              disabled>
            </Input>
            <Input
              v-model="mainStore.stableConfig.height"
              label="Height"
              disabled>
            </Input>
            <Input
              v-model="mainStore.stableConfig.width"
              label="Width"
              disabled>
            </Input>
            <Input
              v-model="mainStore.stableConfig.steps"
              label="Steps"
              disabled>
            </Input>
            <Input
              v-model="mainStore.stableConfig.cfg_scale"
              label="Cfg Scale"
              disabled>
            </Input>
          </div>

        </CardBox>

    </SectionMain>
  </LayoutAuthenticated>
</template>