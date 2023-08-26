<!--
 * @Author: SpenserCai
 * @Date: 2023-08-17 18:23:21
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-27 00:37:58
 * @Description: file content
-->
<div align="center">

<img src="https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/logo.png" width="200" height="200" alt="sd-webui-discord">

# SD-WEBUI-DISCORD
Support For Clustered Stable Diffusion WebUi Discord Bot

</div>

<div align="center">
  <a href="https://raw.githubusercontent.com/SpenserCai/sd-webui-go/main/LICENSE">
    <img src="https://img.shields.io/github/license/SpenserCai/sd-webui-go?color=blueviolet" alt="license">
  </a>
  <img src="https://img.shields.io/badge/Go-1.19+-blue" alt="go">
  <a href="https://github.com/SpenserCai/sd-webui-go/releases">
    <img src="https://img.shields.io/github/v/release/SpenserCai/sd-webui-discord?color=rgb(255%2C0%2C0)&include_prereleases" alt="release">
  </a>
  <a href="https://goreportcard.com/report/github.com/SpenserCai/sd-webui-discord">
    <img src="https://goreportcard.com/badge/github.com/SpenserCai/sd-webui-discord" alt="GoReportCard">
  </a>
  <a href="https://discord.gg/uNJpzEE4sZ">
    <img src="https://discordapp.com/api/guilds/1140177489008807966/widget.png?style=shield"   alt="Discord Server">
  </a>
</div>

SD-WEBUI-DISCORD is a Discord bot developed in Go language for [stable-diffusion-webui](https://github.com/AUTOMATIC1111/stable-diffusion-webui). It utilizes the [sd-webui-go](https://github.com/SpenserCai/sd-webui-go) to invoke the sd-webui API and supports cluster deployment of multiple sd-webui nodes with automatic scheduling and allocation.

## News
### **2023-08-23: Support ControlNet for `txt2img` command**
By using the `controlnet_detect` command to obtain the parameters of ControlNet and filling them into the `controlnet_args` parameter of the `txt2img` command, you can use ControlNet in txt2img.
### **2023-08-22: Support `txt2img` command**
### **2023-08-22: Support `roop` command** 
### **2023-08-20: Support `controlnet_detect` command**

## Usage

The command is still under active development, and there are two ways to experience sd-webui-discord: 
1. Join our Discord Server where you can try out the latest features and contribute by submitting issues and pull requests. 
2. Self-deploy it to have your own sd-webui-discord instance.


### Join Our Discord Server
[![Discord](https://discordapp.com/api/guilds/1140177489008807966/widget.png?style=banner2)](https://discord.gg/uNJpzEE4sZ)

<!--支持收起-->

#### Text to Image
<details>
<summary>See Image</summary>

  ![Demo](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/txt2img_demo.png)

</details>

#### ControlNet
Extension: [sd-webui-controlnet](https://github.com/Mikubill/sd-webui-controlnet)
<details>
<summary>See Image</summary>

  ![Demo1](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/controlnet_1.jpeg)
  ![Demo2](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/controlnet_2.jpeg)

</details>

#### Roop
Extension: [sd-webui-roop](https://github.com/s0md3v/sd-webui-roop)
<details>
<summary>See Image</summary>

  ![Demo](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/roop_demo.jpeg)

</details>

#### Segment Anything With Prompt
Extension: [sd-webui-segment-anything](https://github.com/continue-revolution/sd-webui-segment-anything)
<details>
<summary>See Image</summary>

  ![Demo](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/sam_demo.png)
  ![options](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/sam_options.png)

</details>

#### Deoldify
Extension: [sd-webui-deoldify](https://github.com/SpenserCai/sd-webui-deoldify)
<details>
<summary>See Image</summary>

  ![Demo](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/deoldify_demo.png)
  ![options](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/deoldify_options.png)

</details>

### Self-Deploy
You need to install the following extensions on the SD webui:

[sd-webui-segment-anythin](https://github.com/continue-revolution/sd-webui-segment-anything)

[sd-weubi-deoldify](https://github.com/SpenserCai/sd-webui-deoldify)

[stable-diffusion-webui-rembg](https://github.com/AUTOMATIC1111/stable-diffusion-webui-rembg)

[sd-webui-roop](https://github.com/s0md3v/sd-webui-roop)

[sd-webui-controlnet](https://github.com/Mikubill/sd-webui-controlnet)

***

1.Download the latest release from [here](https://github.com/SpenserCai/sd-webui-discord/releases/latest).

2.Create a bot account on Discord and get the token.

3.Configuration and Startup
```bash
tar -zxvf sd-webui-discord-release-v*.tar.gz // unzip the release package

cd sd-webui-discord-release-v*/release/
```
Edit the `config.json` file and fill in the token and other information.

```json
{
    "sd_webui":{
        "servers":[
            {
                "name":"webui-1",
                "host":"127.0.0.1:7860",
                "max_concurrent":5,
                "max_queue":100,
                "max_vram":"20G"
            }
        ]
    },
    "discord":{
        "token":"<your token here>",
        "server_id":"<your servers id here if empty all servers>"
    }
}
```

If you want set default value with sd-webui
```json
{
    "sd_webui":{
        "servers":[...],
        "default_setting": {
            "cfg_scale": 8,
            "negative_prompt": "bad,text,watermask",
            "height":1024,
            "width":1024
        }
    }
    ...
}
```

Start The Bot
```bash
# if you can't connect discord,you need use proxy and run this command:
# export https_proxy=http://127.0.0.1:8888;export http_proxy=http://127.0.0.1:8888;
./sd-webui-discord
```

## Participating
This is an ongoing project, and if you are interested in contributing, you can join our [Discord Server](https://discord.gg/uNJpzEE4sZ). We welcome any feedback or suggestions, so feel free to submit an issue.

