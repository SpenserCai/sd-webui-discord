<!--
 * @Author: SpenserCai
 * @Date: 2023-10-22 18:37:31
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-11-03 21:27:07
 * @Description: file content
-->
<script setup>
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import UserCard from '@/components/UserCard.vue'
import SectionMain from '@/components/SectionMain.vue'
import CardBox from '@/components/CardBox.vue'
import SectionTitleLine from '@/components/SectionTitleLine.vue'
import { useMainStore } from '@/stores/main'
import { Input,Textarea,Toggle } from 'flowbite-vue'
import { mdiApplicationSettings } from '@mdi/js'
import { watch, ref } from 'vue'
import Cookies from "js-cookie"


const mainStore = useMainStore()

// 从cookie获取nsfw_filter，如果没有则默认为true
const GetNsfwFilter = () => {
  var nsfw_filter = Cookies.get('nsfw_filter')
  if (nsfw_filter == undefined) {
    return true
  } else {
    // 转换为布尔值
    return nsfw_filter == 'true'
  }
}
// value是一个布尔值，true为开启，false为关闭
const SetNsfwFilter = (value) => {
  console.log(value)
  Cookies.set('nsfw_filter', value.toString())
}

// Define a variable to get and set nsfw filter
const nsfwFilterValue = ref(GetNsfwFilter())

const handleNsfwFilterValueChange = (value) => {
  nsfwFilterValue.value = value
  SetNsfwFilter(value)
}

// 监听nsfwFilterValue的变化，如果变化则更新mainStore的nsfwFilter
watch(nsfwFilterValue, (value) => {
  handleNsfwFilterValueChange(value)
})
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <UserCard class="mb-6" />
      <SectionTitleLine :main="false" title="Setting" :icon="mdiApplicationSettings" />
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
            v-model="mainStore.stableConfig.sampler"
            label="Sampler"
            disabled>
          </Input>
          <Input
            v-model="mainStore.stableConfig.CLIP_stop_at_last_layers"
            label="CLIP Skip"
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
          <Textarea v-model="mainStore.stableConfig.negative_prompt" class="col-span-2" rows="4" label="Negative Prompt" />
          <Toggle v-model="nsfwFilterValue" label="NSFW Filter" />
        </div>
      </CardBox>

    </SectionMain>
  </LayoutAuthenticated>
</template>