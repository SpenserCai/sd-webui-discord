<!--
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-04 12:20:27
 * @Description: file content
-->
<script setup>
import SectionFullScreen from '@/components/SectionFullScreen.vue'
import CardBox from '@/components/CardBox.vue'
import BaseButton from '@/components/BaseButton.vue'
import BaseButtons from '@/components/BaseButtons.vue'
import LayoutGuest from '@/layouts/LayoutGuest.vue'
import { useRoute } from 'vue-router'

// 声明一个dict用于存储错误信息，dict的每个item有title和content两个属性
const errorDict = {
  'auth_error': {
    title: 'Auth',
    content: 'Auth error'
  },
  'login_error': {
    title: 'Login',
    content: 'Login error'
  }
}
const route = useRoute()

// 一个方法通过get参数error获取错误信息
const getError = () => {
  let {error} = route.query
  console.log(error)
  // | 分割error第一个是title，第二个是content
  error = error.split('|')
  var message = ""
  if (error.length == 2) {
    message = error[1]
  }

  var errItem = errorDict[error[0]]
  console.log(errItem)
  if (errItem == undefined) {
    errItem = {
      title: 'Unknown',
      content: 'Unknown error'
    }
  } else {
    if (message != "") {
      errItem.content = message
    }
  }
  return errItem
}

</script>

<template>
  <LayoutGuest>
    <SectionFullScreen v-slot="{ cardClass }">
      <CardBox :class="cardClass">
        <div class="space-y-3">
          <h1 class="text-2xl">{{ getError().title }}</h1>

          <p>{{ getError().content }}</p>
        </div>

        <template #footer>
          <BaseButtons>
            <BaseButton label="Done" to="/" color="danger" />
          </BaseButtons>
        </template>
      </CardBox>
    </SectionFullScreen>
  </LayoutGuest>
</template>
