<!--
 * @Author: SpenserCai
 * @Date: 2023-10-06 17:25:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-07 15:55:13
 * @Description: file content
-->
<script setup>
import { ref } from 'vue'
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import SectionMain from '@/components/SectionMain.vue'
import { userhistory } from '@/api/account'
import SectionTitleLine from '@/components/SectionTitleLine.vue'
import { Pagination,Modal,Img,Avatar } from 'flowbite-vue'
import { mdiDrawPen,mdiClose } from '@mdi/js'
// import CardBox from '@/components/CardBox.vue'
// import SectionTitleLine from '@/components/SectionTitleLine.vue'
// import { useMainStore } from '@/stores/main'
// import { mdiApplicationSettings } from '@mdi/js'

// const mainStore = useMainStore()
// 获取用户历史记录方法，参数为页数和每页数量
const isShowImageInfoModal = ref(false)
const currentImageInfo = ref({})
const show = ref(false)
const total = ref(0)
const currentPage = ref(1)
// 当前list
const currentList = ref([])
const getListFunc = (page, pageSize) => {
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
  })
}

const onPageChanged = (page) => {
  console.log(currentPage.value)
  getListFunc(page, 12)
}

const getImage = (index) => {
  let history = currentList.value[index]
  if (history == undefined) {
    return ""
  } else {
    return history.images[0]
  }
}

const showImageInfo = (index) => {
  let history = currentList.value[index]
  currentImageInfo.value = history
  isShowImageInfoModal.value = true
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

const closeImageInfo = () => {
  isShowImageInfoModal.value = false
}

getListFunc(1, 12)
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
        <Modal v-if="isShowImageInfoModal" size="5xl" @close="closeImageInfo">
          <template #body>
            <div class="flex justify-center">
              <Img size="max-w-lg" alt="My gallery" img-class="rounded-lg transition-all duration-300 cursor-pointer filter" :src="getImagesList()[0].src"/>
            </div>
            <div class="content relative mx-auto w-full max-w-5xl rounded-2xl p-4 pt-5 text-white md:p-8 translate-y-0 opacity-100">
              <div class="flex w-full flex-wrap-reverse justify-between">
                <div class="-ml-2 flex items-center justify-between max-md:w-full md:justify-start">
                  <Avatar status="online" rounded :img="currentImageInfo.user_avatar"></Avatar>
                  <p :title="currentImageInfo.user_name" class="ml-4 text-lg font-medium active:text-blue-100 group-hover:underline underline-offset-2 active:underline-offset-4 break-all line-clamp-1 md:text-xl">{{ currentImageInfo.user_name }}</p>
                </div>
              </div>
              <div class="h-4"></div>
              <div class="flex-col items-start grid grid-cols-2 gap-2">
                <div class="flex flex-col items-start ">
                  <SectionTitleLine :main="false" title="Prompt" :icon="mdiDrawPen" />
                  <p class="text-base font-medium text-gray-300">{{ currentImageInfo.options.prompt }}</p>
                </div>
                <div class="flex flex-col items-start ">
                  <SectionTitleLine :main="false" title="Negative Prompt" :icon="mdiClose" />
                  <p class="text-base font-medium text-gray-300">{{ currentImageInfo.options.negative_prompt }}</p>
                </div>
              </div>

            </div>
          </template>
        </Modal>
        <div v-if="show" id="t2i_list" class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <!--循环4次生存4个<div class="grid gap-4">，每个里面有3个div-->
          <div v-for="(number,index) of 4" :key="index" class="grid gap-4">
            <div v-for="(i_number,i_index) of 3" :key="i_index">
              <img class="h-auto max-w-full rounded-lg" :src="getImage(index*3+i_index)" alt="" @click="showImageInfo(index*3+i_index)">
            </div>
          </div>
        </div>
        <div class="lg:text-center my-3">
            <Pagination v-model="currentPage" :total-pages="total/12 + 1" :slice-length="4" @page-changed="onPageChanged"></Pagination>
        </div>
    </SectionMain>
  </LayoutAuthenticated>
</template>