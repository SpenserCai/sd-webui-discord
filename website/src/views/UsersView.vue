<!--
 * @Author: SpenserCai
 * @Date: 2023-10-11 21:36:11
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-12 00:20:57
 * @Description: file content
-->
<script setup>
import SectionMain from '@/components/SectionMain.vue'
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import SectionTitleLine from '@/components/SectionTitleLine.vue'

import { mdiAccountMultiple } from '@mdi/js';
import { onMounted,ref } from 'vue'
import { Table, TableHead, TableBody, TableHeadCell, TableRow, TableCell,Avatar,Pagination,Toggle } from 'flowbite-vue'
import { userlist } from '@/api/system'

const total = ref(0)
const currentPage = ref(1)
const pageSize = 20

const userList = ref([])

const getTotalPage = (total) => {
  // 如果total是12的倍数，返回total/12，否则返回total/12取整+1
  if (total % pageSize == 0) {
    return total / pageSize
  } else {
    return Math.floor(total / pageSize) + 1
  }
}

const onPageChanged = (page) => {
  getUserList(page)
}

const getUserList = (page) => {
    userlist({
      query:{
        id: "",
        name: "",
        only_enable: false
      },
      page_info: {
        page: page,
        page_size: pageSize // 12
      }
    }).then(res => {
      userList.value = res.data.users
      total.value = res.data.page_info.total
      currentPage.value = res.data.page_info.page
    })
}

onMounted(() => {
  getUserList(currentPage.value)
})

</script>

<template>
    <LayoutAuthenticated>
      <SectionMain>
        <SectionTitleLine main title="Users" :icon="mdiAccountMultiple" />
        <Table hoverable>
          <table-head>
            <table-head-cell>Username</table-head-cell>
            <table-head-cell>Enable</table-head-cell>
            <table-head-cell>Roles</table-head-cell>
            <table-head-cell>Created</table-head-cell>
          </table-head>
          <table-body>
            <table-row v-for="(item, index) in userList" :key="index">
              <table-cell>
                <div class="-ml-2 flex items-center justify-between max-md:w-full md:justify-start">
                  <Avatar size="xs" status="online" rounded :img="item.avatar"></Avatar>
                  <p :title="item.name" class="ml-2 text-sm font-medium active:text-blue-100 group-hover:underline underline-offset-2 active:underline-offset-4 break-all line-clamp-1 md:text-sm">{{ item.username }}</p>
                </div>
              </table-cell>
              <table-cell>
                <Toggle v-model="item.enable" :disabled="true" />
              </table-cell>
              <table-cell>{{ item.roles }}</table-cell>
              <table-cell>{{ item.created }}</table-cell>
            </table-row>
          </table-body>
        </Table>
        <div class="lg:text-center my-3">
          <Pagination v-model="currentPage" :total-pages="getTotalPage(total)" :slice-length="4" @page-changed="onPageChanged"></Pagination>
        </div>
      </SectionMain>
    </LayoutAuthenticated>
</template>