/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 10:22:20
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 12:51:09
 * @Description: file content
 */
import {
  mdiAccount,
  mdiLogout,
  mdiGithub,
  mdiPalette,
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
    isDiscordUrl: true,
    icon: mdiPalette,
    label: 'Go to Discord',
    href: '/api/open_discord_server',
    target: '_blank'
    
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
