<!--
 * @Author: SpenserCai
 * @Date: 2023-10-06 17:25:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-27 23:46:24
 * @Description: file content
-->
<script setup>
import { onMounted, ref, watch } from 'vue'
import LayoutAuthenticated from '@/layouts/LayoutAuthenticated.vue'
import SectionMain from '@/components/SectionMain.vue'
import { userhistory,communityhistory } from '@/api/account'
import CardBox from '@/components/CardBox.vue'
import BaseIcon from '@/components/BaseIcon.vue'
import { Modal,Avatar,Button,Spinner } from 'flowbite-vue'
import { mdiDrawPen,mdiCancel,mdiCogOutline,mdiImageArea,mdiImage,mdiAccountGroup, mdiContentCopy } from '@mdi/js'
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

const overCount = ref(0)

// 每页显示多少行
const gridRowCount = ref(12)
const gridColCount = ref(4)
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

const updateToCurrentList = (history,isAppend=false) => {
  // 循环history，把根据history.images创建small_images,同时创建loading
  for (let i = 0; i < history.length; i++) {
    let tmpImages = history[i].images
    let tmpSmallImages = []
    for (let j = 0; j < tmpImages.length; j++) {
      tmpSmallImages.push(GetSmallImageUrl(tmpImages[j],history[i]))
    }
    history[i].small_images = tmpSmallImages
    history[i].loading_image = GetImageLoadUrl(history[i].options.width,history[i].options.height)
  }
  if (isAppend) {
    currentList.value = [...currentList.value, ...history]
  } else {
    currentList.value = history
  }
  let nextTotal = currentPage.value + 1
  if (nextTotal >= getTotalPage(total.value)) {
    overCount.value = gridRowCount.value - Math.floor((total.value - (currentPage.value * gridColCount.value * gridRowCount.value))/gridColCount.value)
    console.log(overCount.value)
  } else {
    overCount.value = 0
  }
}

const getUserHistory = (page, pageSize,isAppend) => {
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
    updateToCurrentList(res.data.history,isAppend)
    currentPage.value = res.data.page_info.page
    show.value = true
  }).finally(() => {
    isLoading.value = false
  })
}

const getCommunityHistory = (page, pageSize,isAppend) => {
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
    updateToCurrentList(res.data.history,isAppend)
    currentPage.value = res.data.page_info.page
    show.value = true
  }).finally(() => {
    isLoading.value = false
  })
}

const getListFunc = (page, pageSize,isAppend=false) => {
  console.log("getListFunc")
  if (isUserHistory()) {
    getUserHistory(page, pageSize,isAppend)
  } else {
    getCommunityHistory(page, pageSize,isAppend)
  }
}

// const onPageChanged = (page) => {
//   getListFunc(page, gridColCount.value * gridRowCount.value)
// }

const GetSmallImageUrl = (url,history) => {
  if (isDiscordImage(url)) {
    // 把cdn.discordapp.com替换为media.discordapp.net
    url = url.replace("cdn.discordapp.com", "media.discordapp.net")
    // 获取长宽
    let tmpImageWidth = history.options.width
    let tmpImageHeight = history.options.height
    // 如果宽高大于等于512，把宽设置为256，高等比例缩放
    if (tmpImageWidth >= 512 || tmpImageHeight >= 512) {
      tmpImageWidth = 384
      tmpImageHeight = Math.floor(tmpImageHeight * 384 / history.options.width)
    }
    let whString = "width=" + tmpImageWidth + "&height=" + tmpImageHeight
    // 如果url中没有?，则在后面加上?,如果结尾的是&，则直接加上whString，否则加上&whString
    if (url.indexOf("?") == -1) {
      url += "?" + whString
    } else if (url[url.length - 1] == "&") {
      url += whString
    } else {
      url += "&" + whString
    }
  }
  return url
}

// 计算图片加载时的base64url
const GetImageLoadUrl = (width,height) => {
  let tW = width
  let tH = height
  let canvas = document.createElement("canvas")
  canvas.width = tW
  canvas.height = tH
  let ctx = canvas.getContext("2d")
  ctx.fillStyle = "#2d3748"
  ctx.fillRect(0, 0, tW, tH)
  return canvas.toDataURL()
}

const getImage = (index,isSmall=false) => {
  // console.log(index)
  let history = currentList.value[index]
  if (history == undefined) {
    return ""
  } else {
    let tmpImage = history.images[0]
    if (isSmall) {
      tmpImage = history.small_images[0]
    }
    return tmpImage
  }
}

const getGalleryImageLoadStartImg = (index) => {
  return currentList.value[index].loading_image
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
  // 计算每一步消耗的时间单位ms,先获取当前时间
  let history = currentList.value[index]
  currentImageInfo.value = history
  isShowImageInfoModal.value = true
  // 延时0.5秒执行，等待Modal显示
  // setTimeout(() => {
  //   // 给image_detail内的第一个div设置style: top: 0
  //   let div = document.getElementById("image_detail").children[1].children[0]
  //   // 如果div的高度大于app高度，设置top为3rem
  //   if (div.clientHeight > window.innerHeight) {
  //     div.style.top = "3rem"
  //   }
  // }, 100)
  
  // 给image_detail内的第一个div设置style: top: 0
  let div = document.getElementById("image_detail").children[1].children[0]
  // 如果div的高度大于app高度，设置top为3rem
  if (div.clientHeight > window.innerHeight) {
    div.style.top = "3rem"
  }
  
}

const getTotalPage = (total) => {
  // 如果total是12的倍数，返回total/12，否则返回total/12取整+1
  if (total % (gridColCount.value * gridRowCount.value) == 0) {
    return total / (gridColCount.value * gridRowCount.value)
  } else {
    return Math.floor(total / (gridColCount.value * gridRowCount.value)) + 1
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
    let tmpWidth = currentImageInfo.value.options.width/2
    let tmpHeight = currentImageInfo.value.options.height/2
    // 如果/2之之后任何一个不是4的倍数，则调整成最接近当前值的4的倍数
    if (tmpWidth % 4 != 0) {
      tmpWidth = Math.floor(tmpWidth / 4) * 4
      // 等比例调整高度(判断新的tmpWidth相比于原来的tmpWidth的比例，然后等比例调整tmpHeight)
      tmpHeight = Math.floor(tmpHeight * (tmpWidth / (currentImageInfo.value.options.width / 2)))
    } else if (tmpHeight % 4 != 0) {
      tmpHeight = Math.floor(tmpHeight / 4) * 4
      // 等比例调整宽度
      tmpWidth = Math.floor(tmpWidth * (tmpHeight / (currentImageInfo.value.options.height / 2)))
    }
    let pixels = decode(blurdata, tmpWidth, tmpHeight);
    base64data = convertUint8ClampedArrayToBase64Image(pixels, tmpWidth, tmpHeight)
    imagesList.push({
      src: images[i],
      alt: "image_" + i,
      hash: base64data
    })
    
  }
  return imagesList
}

const LoadNext = () => {
  if (isLoading.value) {
    console.log("isLoading is true")
    return
  }
  // 判断当前页是否是最后一页，如果是则不执行
  if (currentPage.value == getTotalPage(total.value)) {
    console.log("currentPage is last page")
    return
  }
  // 计算下一页是否最后一页，如果是则计算下一页多余的数量

  getListFunc(currentPage.value + 1, gridColCount.value * gridRowCount.value, true)
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

const galleryImageLoaded = (e) => {
  // 获取id
  let id = e.target.id
  // 获取图片
  let img = document.getElementById(id)
  // if img is null, return
  if (img == null) {
    return
  }
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
    img.hidden = false
  })
}

// const RefreshCurrentPage = () => {
//   getListFunc(currentPage.value, gridColCount.value * gridRowCount.value)
// }

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

  getListFunc(1, gridColCount.value * gridRowCount.value)

})

onscroll = () => {
  const scrollHeight = document.documentElement.scrollHeight
  const scrollTop = document.documentElement.scrollTop
  const clientHeight = document.documentElement.clientHeight
  // 滚动条到最后10%时加载下一页
  if (Math.floor(scrollTop + clientHeight) >= scrollHeight * 0.9) {
    // console.log('快到底了!')
    LoadNext()
  }
}

watch(() => router.currentRoute.value.path,() => {
  // 清空currentList
  show.value = false
  currentList.value = []
  getListFunc(1, gridColCount.value * gridRowCount.value)
})
</script>

<template>
  <LayoutAuthenticated>
    <SectionMain>
      <Modal v-show="isShowImageInfoModal" id="image_detail" size="5xl" @close="closeImageInfo">
        <template #body>
          <div v-if="isShowImageInfoModal" class="content relative mx-auto w-full max-w-5xl rounded-2xl p-4 pt-1 text-white md:p-8 md:pt-1 translate-y-0 opacity-100">
            <CardBox>
              <span class="flex items-center justify-start gap-1 text-xs font-bold uppercase leading-[0] tracking-wide text-slate-400 md:text-sm">
                <BaseIcon :path="mdiImage" class="text-indigo-500" size="16"/>
                Image
              </span>
              <div class="h-2"></div>
              <div class="flex justify-center">
                <img v-lazy="{ src:getImagesList()[0].src, loading: getImagesList()[0].hash, delay: 500}" class="rounded-lg transition-all duration-300 cursor-pointer filter max-w-lg max-h-80" alt="My gallery"/>
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
      <SectionTitleLine v-if="isUserHistory()" main title="Gallery" :icon="mdiImageArea"></SectionTitleLine>
      <SectionTitleLine v-else main title="Community" :icon="mdiAccountGroup"></SectionTitleLine>
      
      <div v-if="show" class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <div v-for="(number,index) of 4" :key="index" class="columns-1">
          <div v-for="(i_number,i_index) of gridRowCount * currentPage - overCount" :key="i_index" class="overflow-hidden rounded-lg mb-4" >
              <img :id="number+'_'+i_number+'_'+'gallery'" v-lazy="{ src: getImage(i_index*4+index,true), loading: getGalleryImageLoadStartImg(i_index*4+index,true), delay: 500}"  crossorigin="anonymous" class="h-auto max-w-full rounded-lg object-cover" alt="" @load="galleryImageLoaded" @click="showImageInfo(i_index*4+index)" >
          </div>
        </div>
      </div>

      <div class="lg:text-center my-3">
        <Button v-if="isLoading" color="default" outline size="xl" >
          <spinner color="blue" />
          <template #suffix>
            <span class="ml-2">Loading More...</span>
          </template>
        </Button>
      </div>
    </SectionMain>
  </LayoutAuthenticated>
</template>