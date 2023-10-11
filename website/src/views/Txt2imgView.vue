<!--
 * @Author: SpenserCai
 * @Date: 2023-10-06 17:25:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-11 02:42:19
 * @Description: file content
-->
<script setup>
import { onMounted, ref, watch } from 'vue'
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import SectionMain from '@/components/SectionMain.vue'
import { userhistory,communityhistory } from '@/api/account'
import CardBox from '@/components/CardBox.vue'
import BaseIcon from '@/components/BaseIcon.vue'
import { Pagination,Modal,Img,Avatar,Button,Spinner } from 'flowbite-vue'
import { mdiDrawPen,mdiCancel,mdiCogOutline,mdiImageArea,mdiImage,mdiAccountGroup, mdiRefresh, mdiContentCopy } from '@mdi/js'
import { useRoute,useRouter } from 'vue-router'
// import CardBox from '@/components/CardBox.vue'
import SectionTitleLine from '@/components/SectionTitleLine.vue'
import NotifyGroup from '@/components/NotifyGroup.vue'
import { notify } from "notiwind"
// import { useMainStore } from '@/stores/main'
// import { mdiApplicationSettings } from '@mdi/js'

// const mainStore = useMainStore()
// 获取用户历史记录方法，参数为页数和每页数量
const route = useRoute()
const isLoading = ref(false)
const router = useRouter()
const isShowImageInfoModal = ref(false)
const currentImageInfo = ref({})
const show = ref(false)
const total = ref(0)
const currentPage = ref(1)
// 每页显示多少行
const gridRowCount = ref(3)
// 当前list
const currentList = ref([])

const isUserHistory = () => {
  let history_type = route.path
  console.log(history_type)
  if (history_type == "/community") {
    return false
  } else {
    return true
  }
}

const getUserHistory = (page, pageSize) => {
  isLoading.value = true
  userhistory({
    query: {
        command: "txt2img"
    },
    page_info: {
        page: page,
        page_size: pageSize // 12
    }
  }).then(res => {
    total.value = res.data.page_info.total
    currentPage.value = res.data.page_info.page
    currentList.value = res.data.history
    show.value = true
  }).finally(() => {
    isLoading.value = false
  })
}

const getCommunityHistory = (page, pageSize) => {
  isLoading.value = true
  communityhistory({
    query: {
        command: "txt2img"
    },
    page_info: {
        page: page,
        page_size: pageSize // 12
    }
  }).then(res => {
    total.value = res.data.page_info.total
    currentPage.value = res.data.page_info.page
    currentList.value = res.data.history
    show.value = true
  }).finally(() => {
    isLoading.value = false
  })
}

const getListFunc = (page, pageSize) => {
  if (isUserHistory()) {
    getUserHistory(page, pageSize)
  } else {
    getCommunityHistory(page, pageSize)
  }
}

const onPageChanged = (page) => {
  console.log(currentPage.value)
  getListFunc(page, 4 * gridRowCount.value)
}

const getImage = (index) => {
  let history = currentList.value[index]
  if (history == undefined) {
    return ""
  } else {
    return history.images[0]
  }
}

const getCurrentImageVae = () => {
  let vae = currentImageInfo.value.options.override_settings.sd_vae
  if (vae == undefined) {
    return "Automatic"
  } else {
    return vae.split(".")[0]
  }
}

const getCurrentImageModel = () => {
  let model = currentImageInfo.value.options.override_settings.sd_model_checkpoint
  if (model == undefined) {
    return "Default"
  } else {
    return model.split(".")[0]
  }
}

const showImageInfo = (index) => {
  let history = currentList.value[index]
  currentImageInfo.value = history
  isShowImageInfoModal.value = true
  // 延时0.5秒执行，等待Modal显示
  setTimeout(() => {
    // 给image_detail内的第一个div设置style: top: 0
    let div = document.getElementById("image_detail").children[1].children[0]
    // 如果div的高度大于app高度，设置top为3rem
    if (div.clientHeight > window.innerHeight) {
      div.style.top = "3rem"
    }
  }, 100)
}

const getTotalPage = (total) => {
  // 如果total是12的倍数，返回total/12，否则返回total/12取整+1
  if (total % 12 == 0) {
    return total / 12
  } else {
    return Math.floor(total / 12) + 1
  }
}

const getImagesList = () => {
  let images = currentImageInfo.value.images
  let imagesList = []
  for (let i = 0; i < images.length; i++) {
    imagesList.push({
      src: images[i],
      alt: "image_" + i
    })
  }
  return imagesList
}

const copyCommand = () => {
  // 获取当前Image的StableConfig
  let stableConfig = currentImageInfo.value.options
  let mainCmd = currentImageInfo.value.command
  let fullCmd = "/" + mainCmd + " "
  // 循环stableConfig，拼接fullCmd，格式是key:[space]value[space],如果value是json则递归
  for (let key in stableConfig) {
    let value = stableConfig[key]
    if (typeof(value) == "object") {
      for (let subKey in value) {
        let subValue = value[subKey]
        if (subKey == "sd_model_checkpoint") {
          subKey = "checkpoint"
        }
        fullCmd += subKey + ": " + subValue + " "
      }
    } else {
      if (key == "sampler_index") {
        key = "sampler"
      } 
      fullCmd += key + ": " + value + " "
    }
  }
  // 把fullCmd复制到剪贴板
  navigator.clipboard.writeText(fullCmd)
  notify({
    title: "Success",
    text: "Command Copied to clipboard",
    type: "success",
    group: "t2i_image_info",
  }, 5000)

}

const closeImageInfo = () => {
  isShowImageInfoModal.value = false
}

onMounted(() => {
  getListFunc(1, 4 * gridRowCount.value)
})

watch(() => router.currentRoute.value.path,() => {
  getListFunc(1, 4 * gridRowCount.value)
})
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <Modal v-if="isShowImageInfoModal" id="image_detail" size="5xl" @close="closeImageInfo">
        <template #body>
          <div class="content relative mx-auto w-full max-w-5xl rounded-2xl p-4 pt-1 text-white md:p-8 md:pt-1 translate-y-0 opacity-100">
            <CardBox>
              <span class="flex items-center justify-start gap-1 text-xs font-bold uppercase leading-[0] tracking-wide text-slate-400 md:text-sm">
                <BaseIcon :path="mdiImage" class="text-indigo-500" size="16"/>
                Image
              </span>
              <div class="h-2"></div>
              <div class="flex justify-center">
                <Img size="max-w-lg max-h-80" alt="My gallery" img-class="rounded-lg transition-all duration-300 cursor-pointer filter" :src="getImagesList()[0].src"/>
              </div>
              <div class="h-2"></div>
              <div class="flex w-full flex-wrap-reverse justify-between">
                <div class="-ml-2 flex items-center justify-between max-md:w-full md:justify-start">
                  <Avatar status="online" rounded :img="currentImageInfo.user_avatar"></Avatar>
                  <p :title="currentImageInfo.user_name" class="ml-4 text-lg font-medium active:text-blue-100 group-hover:underline underline-offset-2 active:underline-offset-4 break-all line-clamp-1 md:text-xl">{{ currentImageInfo.user_name }}</p>
                </div>
              </div>
            </CardBox>
            <div class="h-4"></div>
            <div class="flex-col items-start grid grid-cols-2 gap-2">
              <CardBox class="flex flex-col items-start h-full">
                <span class="flex items-center justify-start gap-1 text-xs font-bold uppercase leading-[0] tracking-wide text-slate-400 md:text-sm">
                  <BaseIcon :path="mdiDrawPen" class="text-emerald-600" size="16"/>
                  Prompt
                </span>
                <div class="h-2"></div>
                <div class="max-h-32 overflow-y-auto aside-scrollbars dark:aside-scrollbars-[slate]">
                <p class="text-base font-medium text-gray-300 break-all">{{ currentImageInfo.options.prompt }}</p>
                </div>
              </CardBox>
              <CardBox class="flex flex-col items-start h-full">
                <span class="flex items-center justify-start gap-1 text-xs font-bold uppercase leading-[0] tracking-wide text-slate-400 md:text-sm">
                  <BaseIcon :path="mdiCancel" class="text-red-600" size="16"/>
                  Negative Prompt
                </span>
                <div class="h-2"></div>
                <div class="max-h-32 overflow-y-auto aside-scrollbars dark:aside-scrollbars-[slate]">
                <p class="text-base font-medium text-gray-300 break-all">{{ currentImageInfo.options.negative_prompt }}</p>
                </div>
              </CardBox>
            </div>
            <div class="h-4"></div>
            <div class="flex-col items-start">
              <CardBox>
                <span class="flex items-center justify-start gap-1 text-xs font-bold uppercase leading-[0] tracking-wide text-slate-400 md:text-sm">
                  <BaseIcon :path="mdiCogOutline" class="text-indigo-500" size="16"/>
                  Info
                </span>
                <div class="h-2"></div>
                <div class="grid grid-cols-3 gap-3">
                  <div title="Model" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Model</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ getCurrentImageModel() }}</p>
                  </div>
                  <div title="Vae" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Vae</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ getCurrentImageVae() }}</p>
                  </div>
                  <div title="Sampler" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Sampler</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ currentImageInfo.options.sampler_index }}</p>
                  </div>
                  <div title="Size" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Size</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ currentImageInfo.options.height }}<span class="mx-1 opacity-50">x</span>{{ currentImageInfo.options.width }}</p>
                  </div>
                  <div title="Steps" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Steps</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ currentImageInfo.options.steps }}</p>
                  </div>
                  <div title="Cfg Scale" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Cfg Scale</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ currentImageInfo.options.cfg_scale }}</p>
                  </div>
                  <div title="Seed" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Seed</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ currentImageInfo.options.seed }}</p>
                  </div>
                  <div title="Created" class="flex flex-col items-start justify-start leading-tight">
                    <h5 class="text-xs font-medium tracking-wider text-slate-500">Created</h5>
                    <p class="whitespace-nowrap max-md:text-sm">{{ currentImageInfo.created }}</p>
                  </div>
                  <div>
                    <div class="h-1"></div>
                    <Button size="xs" gradient="teal-lime" @click="copyCommand()">
                      Copy Command
                      <template #suffix>
                        <BaseIcon :path="mdiContentCopy" />
                      </template>
                    </Button>
                  </div>
                </div>
              </CardBox>
            </div>
          </div>
          <NotifyGroup group="t2i_image_info" />
        </template>
      </Modal>
      <SectionTitleLine v-if="isUserHistory()" main title="Gallery" :icon="mdiImageArea">
        <Button size="xs" gradient="purple-blue" outline square @click="getListFunc(1, 4 * gridRowCount)">
          <spinner v-show="isLoading" size="6" />
          <BaseIcon :path="mdiRefresh" />
        </Button>
      </SectionTitleLine>
      <SectionTitleLine v-else main title="Community" :icon="mdiAccountGroup">
        <Button size="xs" gradient="purple-blue" outline square @click="getListFunc(1, 4 * gridRowCount)">
          <spinner v-show="isLoading" size="6" />
          <BaseIcon :path="mdiRefresh" />
        </Button>
      </SectionTitleLine>
      
      <div v-if="show" id="t2i_list" class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <!--循环4次生存4个<div class="grid gap-4">，每个里面有3个div-->
        <div v-for="(number,index) of 4" :key="index" class="grid gap-4">
          <div v-for="(i_number,i_index) of gridRowCount" :key="i_index">
            <img class="h-auto max-w-full rounded-lg" :src="getImage(i_index*4+index)" alt="" @click="showImageInfo(i_index*4+index)">
          </div>
        </div>
      </div>
      <div class="lg:text-center my-3">
          <Pagination v-model="currentPage" :total-pages="getTotalPage(total)" :slice-length="4" @page-changed="onPageChanged"></Pagination>
      </div>
    </SectionMain>
  </LayoutAuthenticated>
</template>