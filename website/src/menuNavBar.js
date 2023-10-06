/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-07 00:18:25
 * @Description: file content
 */
import {
  mdiAccount,
  mdiLogout,
  mdiGithub,
} from '@mdi/js'

export default [
  {
    isCurrentUser: true,
    menu: [
      {
        icon: mdiAccount,
        label: 'Profile',
        to: '/app'
      },
      {
        isDivider: true
      },
      {
        icon: mdiLogout,
        label: 'Log Out',
        isLogout: true
      }
    ]
  },
  {
    icon: mdiGithub,
    label: 'GitHub',
    isDesktopNoLabel: true,
    href: 'https://github.com/SpenserCai/sd-webui-discord',
    target: '_blank'
  },
  {
    icon: mdiLogout,
    label: 'Log out',
    isDesktopNoLabel: true,
    isLogout: true
  }
]
