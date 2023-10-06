<!--
 * @Author: SpenserCai
 * @Date: 2023-10-06 17:25:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-06 21:59:28
 * @Description: file content
-->
<script setup>
import { ref } from 'vue'
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import SectionMain from '@/components/SectionMain.vue'
import { userhistory } from '@/api/account'
import { Pagination } from 'flowbite-vue'
// import CardBox from '@/components/CardBox.vue'
// import SectionTitleLine from '@/components/SectionTitleLine.vue'
// import { useMainStore } from '@/stores/main'
// import { mdiApplicationSettings } from '@mdi/js'

// const mainStore = useMainStore()
// 获取用户历史记录方法，参数为页数和每页数量
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

getListFunc(1, 12)
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
        <div v-if="show" id="t2i_list" class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <!--循环4次生存4个<div class="grid gap-4">，每个里面有3个div-->
          <div v-for="(number,index) of 4" :key="index" class="grid gap-4">
            <div v-for="(i_number,i_index) of 3" :key="i_index">
                <img class="h-auto max-w-full rounded-lg" :src="getImage(index*3+i_index)" alt="">
            </div>
          </div>
        </div>
        <div class="lg:text-center my-3">
            <Pagination v-model="currentPage"  :total-pages="total/12 + 1" :slice-length="4" @page-changed="onPageChanged"></Pagination>
        </div>

    </SectionMain>
  </LayoutAuthenticated>
</template>