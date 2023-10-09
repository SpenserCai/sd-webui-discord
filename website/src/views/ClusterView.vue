<!--
 * @Author: SpenserCai
 * @Date: 2023-10-09 16:47:39
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 18:20:42
 * @Description: file content
-->
<script setup>
import SectionTitleLine from '@/components/SectionTitleLine.vue'
import SectionMain from '@/components/SectionMain.vue'
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import CardBox from '@/components/CardBox.vue'
import { mdiServer } from '@mdi/js'
import BaseIcon from '@/components/BaseIcon.vue'
import { ref } from 'vue'
import { cluster } from '@/api/system'

const clusterNodes = ref([])

const getListFunc = () => {
    cluster().then(res => {
        clusterNodes.value = res.data.cluster
    })
}

getListFunc()

</script>

<template>
    <LayoutAuthenticated>
      <SectionMain>
        <SectionTitleLine main title="Cluster" :icon="mdiServer" />
        <div class="flex-col items-start grid grid-cols-2 gap-2">
            <CardBox v-for="(item, index) in clusterNodes" :key="index">
                <div class="grid grid-cols-3 gap-3">
                    <BaseIcon :path="mdiServer" size="72" w="" h="h-32" class="col-span-1 justify-start"/>
                    <div class="col-span-2">
                        <div class="grid grid-cols-2 gap-2">
                            <div title="Name" class="flex flex-col items-start justify-start leading-tight">
                              <h5 class="text-sm font-medium tracking-wider text-slate-500">Name</h5>
                              <p class="whitespace-nowrap max-md:text-xs">{{ item.name }}</p>
                            </div>
                            <div title="Host" class="flex flex-col items-start justify-start leading-tight">
                              <h5 class="text-sm font-medium tracking-wider text-slate-500">Host</h5>
                              <p class="whitespace-nowrap max-md:text-xs">{{ item.host }}</p>
                            </div>
                            <div title="Running" class="flex flex-col items-start justify-start leading-tight">
                              <h5 class="text-sm font-medium tracking-wider text-slate-500">Running</h5>
                              <p class="whitespace-nowrap max-md:text-xs">{{ item.running }}</p>
                            </div>
                            <div title="Pending" class="flex flex-col items-start justify-start leading-tight">
                              <h5 class="text-sm font-medium tracking-wider text-slate-500">Pending</h5>
                              <p class="whitespace-nowrap max-md:text-xs">{{ item.pending }}</p>
                            </div>
                            <div title="MaxConcurrent" class="flex flex-col items-start justify-start leading-tight col-span-1">
                              <h5 class="text-sm font-medium tracking-wider text-slate-500">Max Concurrent</h5>
                              <p class="whitespace-nowrap max-md:text-xs">{{ item.max_concurrent }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </CardBox>
        </div>

      </SectionMain>
    </LayoutAuthenticated>
</template>