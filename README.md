<!--
 * @Author: SpenserCai
 * @Date: 2023-08-17 18:23:21
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-15 12:14:22
 * @Description: file content
-->
<div align="center">

<img src="https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/logo.png" width="200" height="200" alt="sd-webui-discord">

# SD-WEBUI-DISCORD
Support For Clustered Stable Diffusion WebUi Discord Bot

[Mult-Language README](https://github.com/SpenserCai/sd-webui-discord/tree/main/docs)
</div>

<div align="center">
  <a href="https://raw.githubusercontent.com/SpenserCai/sd-webui-go/main/LICENSE">
    <img src="https://img.shields.io/github/license/SpenserCai/sd-webui-go?color=blueviolet" alt="license">
  </a>
  <img src="https://img.shields.io/badge/Go-1.19+-blue" alt="go">
  <a href="https://github.com/SpenserCai/sd-webui-discord/releases">
    <img src="https://img.shields.io/github/v/release/SpenserCai/sd-webui-discord?color=rgb(255%2C0%2C0)&include_prereleases" alt="release">
  </a>
  <a href="https://goreportcard.com/report/github.com/SpenserCai/sd-webui-discord">
    <img src="https://goreportcard.com/badge/github.com/SpenserCai/sd-webui-discord" alt="GoReportCard">
  </a>
  <a href="https://discord.gg/uNJpzEE4sZ">
    <img src="https://discordapp.com/api/guilds/1140177489008807966/widget.png?style=shield"   alt="Discord Server">
  </a>
  <a href="https://qun.qq.com/qqweb/qunpro/share?_wv=3&_wwv=128&appChannel=share&inviteCode=21gYdX0DSw2&businessType=7&from=181074&biz=ka">
    <img src="https://img.shields.io/badge/QQ%E9%A2%91%E9%81%93-SD%20WEBUI%20DISCORD-5492ff?style=flat-square" alt="QQ Channel">
  </a>
</div>

## Introduction
SD-WEBUI-DISCORD is a Discord bot developed in Go language for [stable-diffusion-webui](https://github.com/AUTOMATIC1111/stable-diffusion-webui). It utilizes the [sd-webui-go](https://github.com/SpenserCai/sd-webui-go) to invoke the sd-webui API and supports cluster deployment of multiple sd-webui nodes with automatic scheduling and allocation.
At the same time, there is also the [sd-webui-discord-ex](https://github.com/SpenserCai/sd-webui-discord-ex), which is an extension on the stable-diffusion-webui that you can install and use directly. It will automatically update every time you restart SD webui.

## Screenshots

![First](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/first_page_new.png)
![Second](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/second_page_new.png)

## News
 - `2023-10-12 `: Support `Website` command
   <details>
     <summary>See Image</summary>
       
      ![example](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/website_gallery.png)
      ![example](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/website_community.png)
      ![example](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/website_image_detail.png)

   </details>
 - `2023-09-24 `: Support multi image generate in `txt2img` command**
 - `2023-09-23 `: Support `Retry` and `Delete` in `txt2img` command
 - `2023-09-22 ` 
     - Support `setting_ui` command
     - Better `txt2img` response ui,thanks for [venetanji](https://github.com/venetanji) support that! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
     - Optimize command add when bot start,thanks for [venetanji](https://github.com/venetanji) support that! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
 - `2023-09-10 `: Support local language
 - `2023-09-05 `: Support User Center on Windows
 - `2023-09-04 `: Support Image to Image
 - `2023-08-31 `: Support User Center
 - `2023-08-27 `:
     - Support txt2img choice model checkpoint
     - Support upload image with attachment: `deoldify` `png_info` `roop_image` commands
     <details>
     <summary>See Image</summary>

      ![example](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/support_attechment.png)
     
     </details>
  
 - `2023-08-26 `: Support Setting default sd-webui options
 - `2023-08-23 `: Support ControlNet for `txt2img` command
     - By using the `controlnet_detect` command to obtain the parameters of ControlNet and filling them into the `controlnet_args` parameter of the `txt2img` command, you can use ControlNet in txt2img.
 - `2023-08-22 `: Support `txt2img` command
 - `2023-08-22 `: Support `roop` command
 - `2023-08-20 `: Support `controlnet_detect` command  

## Features
- Website
    - User gallery.
    - Community gallery.
    - User setting view.
    - User list(Admin).
    - Cluster node list(Admin).
- Local language support.
- Global sd-webui default options.
- Supports multi-node (sd-webui) deployment, distributed cluster queue with automatic scheduling.
- User Center
    - Set user options with ui. 
    - User Center can be freely enabled.
    - Customizable database types, currently supporting MySQL and SQLite.
    - Users can set their own default options. If not specified when generating images, the user's configured options will be used by default, such as image dimensions, cfg, steps, etc. For more details, please refer to the User Center (wiki).
    - Supports user registration.
- ControlNet Preview
    - Supports specifying module and model through selection, eliminating the need for manual input.
    - Allows previewing preprocessing effects and obtaining args simultaneously (used for user-generated txt2img).
- Text to Image
    - Support SDXL's refiner!
    - Supports user-defined default options from the User Center.
    - Enables specifying model, sampler, and other optional parameters through selection, eliminating the need for manual input. 
    - Supports using parameters from "ControlNet Preview" directly. 
- Image to Image
    - Allows uploading images through an image control, rather than using URLs.
    - Supports all operations of img2img in sd-webui!
- Roop Face Swap
    - Allows uploading images through an image control, rather than using URLs.
    - Supports specifying source and target through selection, eliminating the need for manual input.
    - Supports custom face rendering algorithms (GFPGAN, CodeFormer), as well as weights.
- Deoldify Colorization
    - Allows uploading images through an image control, rather than using URLs.
    - Currently the best photo colorization model.
- Segment-Anything
    - Allows uploading images through an image control, rather than using URLs.
    - Supports image segmentation based on natural language descriptions (DION+SAM).
    - Enables specifying DION and SAM models through selection, eliminating the need for manual input.
- Background Removal
    - Allows uploading images through an image control, rather than using URLs.
    - Supports commonly used background removal algorithms.
    - Supports returning the mask.
- Extra Single
    - Allows uploading images through an image control, rather than using URLs.
    - Supports facial repair.
    - Supports super-resolution, with models available for direct selection without manual input.
- Png Info
    - Allows uploading images through an image control, rather than using URLs.
    - Supports retrieving parameters of images generated by sd-webui.

## Usage

The command is still under active development, and there are two ways to experience sd-webui-discord: 
1. Join our Discord Server where you can try out the latest features and contribute by submitting issues and pull requests. 
2. Self-deploy it to have your own sd-webui-discord instance.

[SD-WEBUI-DISCORD Dashboard](https://aigc.ngrok.io/)

[![Discord](https://invidget.switchblade.xyz/uNJpzEE4sZ)](https://discord.gg/uNJpzEE4sZ)


## Documentation
Detailed tutorial reference [Wiki](https://github.com/SpenserCai/sd-webui-discord/wiki)


## Participating
This is an ongoing project, and if you are interested in contributing, you can join our [Discord Server](https://discord.gg/uNJpzEE4sZ). We welcome any feedback or suggestions, so feel free to submit an issue.

## Used By

### WAN Show Bingo!

[![Discord](https://invidget.switchblade.xyz/pWS5mw7jFz)](https://discord.gg/pWS5mw7jFz)
