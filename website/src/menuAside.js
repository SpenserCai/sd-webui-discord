/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 18:23:40
 * @Description: file content
 */
import {
  // mdiAccountCircle,
  // mdiMonitor,
  // mdiLock,
  // mdiAlertCircle,
  // mdiSquareEditOutline,
  // mdiTable,
  // mdiViewList,
  // mdiTelevisionGuide,
  // mdiResponsive,
  // mdiPalette,
  mdiHome,
  mdiImageArea,
  mdiShieldCrown
} from '@mdi/js'

export default [
  {
    to: '/app',
    icon: mdiHome,
    label: 'Home'
  },
  {
    label: 'Gallery',
    icon: mdiImageArea,
    menu:[
      {
        to: '/txt2img',
        label: 'Txt2img'
      }
    ]
  },
  {
    isDivider: true
  },
  {
    label: 'Admin',
    icon: mdiShieldCrown,
    menu:[
      {
        label: 'Cluster',
        to: '/cluster'
      },
      {
        label: 'Users',
        // to: '/users'
      }
    ],
    roles: ['admin']
  }
  // {
  //   to: '/dashboard',
  //   icon: mdiMonitor,
  //   label: 'Dashboard'
  // },
  // {
  //   to: '/tables',
  //   label: 'Tables',
  //   icon: mdiTable
  // },
  // {
  //   to: '/forms',
  //   label: 'Forms',
  //   icon: mdiSquareEditOutline
  // },
  // {
  //   to: '/ui',
  //   label: 'UI',
  //   icon: mdiTelevisionGuide
  // },
  // {
  //   to: '/responsive',
  //   label: 'Responsive',
  //   icon: mdiResponsive
  // },
  // {
  //   to: '/',
  //   label: 'Styles',
  //   icon: mdiPalette
  // },
  // {
  //   to: '/profile',
  //   label: 'Profile',
  //   icon: mdiAccountCircle
  // },
  // {
  //   to: '/login',
  //   label: 'Login',
  //   icon: mdiLock
  // },
  // {
  //   to: '/error',
  //   label: 'Error',
  //   icon: mdiAlertCircle
  // },
  // {
  //   label: 'Dropdown',
  //   icon: mdiViewList,
  //   menu: [
  //     {
  //       label: 'Item One'
  //     },
  //     {
  //       label: 'Item Two'
  //     }
  //   ]
  // }
]
