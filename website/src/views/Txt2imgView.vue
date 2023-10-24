<!--
 * @Author: SpenserCai
 * @Date: 2023-10-06 17:25:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-24 09:08:32
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
import { decode } from "blurhash"
// import * as tf from '@tensorflow/tfjs'
import * as nsfwjs from 'nsfwjs'
// import { useMainStore } from '@/stores/main'
// import { mdiApplicationSettings } from '@mdi/js'

// const mainStore = useMainStore()
// 获取用户历史记录方法，参数为页数和每页数量
const route = useRoute()
const isLoading = ref(true)
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

const model = ref(null)

const isUserHistory = () => {
  let history_type = route.path
  console.log(history_type)
  if (history_type == "/community") {
    return false
  } else {
    return true
  }
}

function convertUint8ClampedArrayToBase64Image(uint8Array, width, height) {
  const imageData = new ImageData(uint8Array, width, height);
  const canvas = document.createElement('canvas');
  const context = canvas.getContext('2d');
  canvas.width = width;
  canvas.height = height;
  context.putImageData(imageData, 0, 0);
  return canvas.toDataURL();
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

// const onPageChanged = (page) => {
//   getListFunc(page, 4 * gridRowCount.value)
// }

const getImage = (index,isSmall=false) => {
  let history = currentList.value[index]
  if (history == undefined) {
    return ""
  } else {
    let tmpImage = history.images[0]
    if (isSmall && isDiscordImage(tmpImage)) {
      // 把cdn.discordapp.com替换为media.discordapp.net
      tmpImage = tmpImage.replace("cdn.discordapp.com", "media.discordapp.net")
      // 获取长宽
      let tmpImageWidth = history.options.width
      let tmpImageHeight = history.options.height
      // 如果宽高大于等于1024，把宽设置为512，高等比例缩放
      if (tmpImageWidth >= 1024 || tmpImageHeight >= 1024) {
        tmpImageWidth = 512
        tmpImageHeight = Math.floor(tmpImageHeight * 512 / history.options.width)
      }
      let whString = "width=" + tmpImageWidth + "&height=" + tmpImageHeight
      // 如果url中没有?，则在后面加上?,如果结尾的是&，则直接加上whString，否则加上&whString
      if (tmpImage.indexOf("?") == -1) {
        tmpImage += "?" + whString
      } else if (tmpImage[tmpImage.length - 1] == "&") {
        tmpImage += whString
      } else {
        tmpImage += "&" + whString
      }
    }
    return tmpImage
  }
}

const getGalleryImageLoadStartImg = (index) => {
  let history = currentList.value[index]
  if (history == undefined) {
    return ""
  } else {
    let tW = history.options.width
    let tH = history.options.height
    // 创建一个宽高为tW,tH的灰色canvas
    let canvas = document.createElement("canvas")
    canvas.width = tW
    canvas.height = tH
    let ctx = canvas.getContext("2d")
    ctx.fillStyle = "#2d3748"
    ctx.fillRect(0, 0, tW, tH)
    return canvas.toDataURL()
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

const isDiscordImage = (url) => {
  if (url.indexOf("discord") != -1) {
    return true
  } else {
    return false
  }
}

const getImagesList = () => {
  let images = currentImageInfo.value.images
  let imagesList = []
  for (let i = 0; i < images.length; i++) {
    let base64data = ""
    let blurdata = currentImageInfo.value.images_blurhash[i]
    if (blurdata == undefined) {
      blurdata = "KED+rLozE4~UakohE4IW%3"
    }  
    
    let pixels = decode(blurdata, currentImageInfo.value.options.width/2, currentImageInfo.value.options.height/2);
    base64data = convertUint8ClampedArrayToBase64Image(pixels, currentImageInfo.value.options.width/2, currentImageInfo.value.options.height/2)
    imagesList.push({
      src: images[i],
      alt: "image_" + i,
      hash: base64data
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
        if (subKey == "controlnet") {
          subKey = "controlnet_args"
          // 把数组转换为字符串逗号连接
          let tmpValue = subValue.args
          let returnValue = ""
          for (let i = 0; i < tmpValue.length; i++) {
            returnValue += JSON.stringify(tmpValue[i]) + ","
          }
          subValue = returnValue.substring(0, returnValue.length - 1)
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

// image detail 的 onload事件
const imageDetailOnload = () => {
  document.getElementById('img_info_loaded').hidden=false
  document.getElementById('img_info_loading').hidden=true
}

const galleryImageLoaded = (e) => {
  // 获取id
  let id = e.target.id
  // 获取图片
  let img = document.getElementById(id)
  let imgLoading = document.getElementById(id + "_loading")
  model.value.classify(img).then( predictions => {
    // 判断predictions是否成功
    switch (predictions[0].className) {
      case 'Hentai':
      case 'Porn':
      case 'Sexy':
        if (predictions[0].probability >= 0.51) {
          img.classList.add("blur-2xl")
        } else {
          img.classList.remove("blur-2xl")
        }
        break
      default:
        img.classList.remove("blur-2xl")
        break
    }
    imgLoading.hidden = true
    img.hidden = false
  })
}

const RefreshCurrentPage = () => {
  getListFunc(currentPage.value, 4 * gridRowCount.value)
}

const closeImageInfo = () => {
  isShowImageInfoModal.value = false
}

onMounted(async () => {
  // 先尝试从indexdb中获取model，如果失败则从/model/中加载，并保存到indexdb中
  try{
    const locationLoaded = await Object.freeze(nsfwjs.load('indexeddb://nsfwjs_model',{size:299}))
    model.value = Object.freeze(locationLoaded)
  } catch (e) {
    console.log("load from /model/")
    await nsfwjs.load('/model/',{size: 299}).then((m) => {
      // https://stackoverflow.com/questions/67815952/vue3-app-with-tensorflowjs-throws-typeerror-cannot-read-property-backend-of-u
      model.value = Object.freeze(m)
      m.model.save('indexeddb://nsfwjs_model')
    })
  }
  getListFunc(1, 4 * gridRowCount.value)
})

// 在currentPage变化时，获取list
watch(currentPage, (newVal, oldVal) => {
  if (newVal == oldVal) {
    console.log("currentPage not changed")
  } else {
    getListFunc(newVal, 4 * gridRowCount.value)
  }
})

watch(() => router.currentRoute.value.path,() => {
  console.log("path changed")
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
                <Img id="img_info_loading" size="max-w-lg max-h-80" alt="My gallery" img-class="rounded-lg transition-all duration-300 cursor-pointer filter" :src="getImagesList()[0].hash"/>
                <Img id="img_info_loaded" hidden="hidden" size="max-w-lg max-h-80" alt="My gallery" img-class="rounded-lg transition-all duration-300 cursor-pointer filter" :src="getImagesList()[0].src" @load="imageDetailOnload"/>
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
        <Button size="xs" gradient="purple-blue" outline square @click="RefreshCurrentPage()">
          <spinner v-show="isLoading" size="6" />
          <BaseIcon :path="mdiRefresh" />
        </Button>
      </SectionTitleLine>
      <SectionTitleLine v-else main title="Community" :icon="mdiAccountGroup">
        <Button size="xs" gradient="purple-blue" outline square @click="RefreshCurrentPage()">
          <spinner v-show="isLoading" size="6" />
          <BaseIcon :path="mdiRefresh" />
        </Button>
      </SectionTitleLine>
      
      <div v-if="show" id="t2i_list" class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <!--循环4次生存4个<div class="grid gap-4">，每个里面有3个div-->
        <div v-for="(number,index) of 4" :key="index" class="grid gap-4">
          <div v-for="(i_number,i_index) of gridRowCount" :key="i_index" class="overflow-hidden rounded-lg">
              <img :id="number+'_'+i_number+'_'+'gallery_loading'" crossorigin="anonymous" class="h-auto animate-pulse max-w-full rounded-lg object-cover" :src="getGalleryImageLoadStartImg(i_index*4+index,true)" alt="" >
              <img :id="number+'_'+i_number+'_'+'gallery'" hidden="hidden" crossorigin="anonymous" class="h-auto max-w-full rounded-lg object-cover" :src="getImage(i_index*4+index,true)" alt="" @load="galleryImageLoaded" @click="showImageInfo(i_index*4+index)">
          </div>
        </div>
      </div>
      <div class="lg:text-center my-3">
          <Pagination v-model="currentPage" :total-pages="getTotalPage(total)" :slice-length="4"></Pagination>
      </div>
    </SectionMain>
  </LayoutAuthenticated>
</template>