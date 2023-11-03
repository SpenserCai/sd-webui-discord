/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-11-03 13:21:33
 * @Description: file content
 */
import { createRouter,createWebHashHistory } from 'vue-router'
import FirstPage from '@/views/FirstView.vue'
import Home from '@/views/HomeView.vue'

const routes = [
  {
    meta: {
      title: 'Welcome'
    },
    path: '/',
    name: 'first',
    component: FirstPage
  },
  {
    meta:{
      title: 'Home'
    },
    path: '/app',
    name: 'app',
    component: () => import('@/views/AppView.vue')
  },
  {
    meta: {
      title: 'Txt2img'
    },
    path: '/txt2img',
    name: 'txt2img',
    component: () => import('@/views/Txt2imgView.vue')
  },
  {
    meta: {
      title: 'Community'
    },
    path: '/community',
    name: 'community',
    component: () => import('@/views/Txt2imgView.vue')
  },
  {
    meta: {
      title: 'Cluster'
    },
    path: '/cluster',
    name: 'cluster',
    component: () => import('@/views/ClusterView.vue')
  },
  {
    meta: {
      title: 'Users'
    },
    path: '/users',
    name: 'users',
    component: () => import('@/views/UsersView.vue')
  },
  {
    // Document title tag
    // We combine it with defaultDocumentTitle set in `src/main.js` on router.afterEach hook
    meta: {
      title: 'Dashboard'
    },
    path: '/dashboard',
    name: 'dashboard',
    component: Home
  },
  {
    meta: {
      title: 'Tables'
    },
    path: '/tables',
    name: 'tables',
    component: () => import('@/views/TablesView.vue')
  },
  {
    meta: {
      title: 'Forms'
    },
    path: '/forms',
    name: 'forms',
    component: () => import('@/views/FormsView.vue')
  },
  {
    meta: {
      title: 'Profile'
    },
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/ProfileView.vue')
  },
  {
    meta: {
      title: 'Ui'
    },
    path: '/ui',
    name: 'ui',
    component: () => import('@/views/UiView.vue')
  },
  {
    meta: {
      title: 'Responsive layout'
    },
    path: '/responsive',
    name: 'responsive',
    component: () => import('@/views/ResponsiveView.vue')
  },
  {
    meta: {
      title: 'Login'
    },
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue')
  },
  {
    meta: {
      title: 'Error'
    },
    path: '/error',
    name: 'error',
    component: () => import('@/views/ErrorView.vue')
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  // history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { top: 0 }
  }
})

export default router
