<!--
 * @Author: SpenserCai
 * @Date: 2023-08-17 18:23:21
 * @version: 
 * @LastEditors: SpenserCai
 * @languageThai: UIXROV
 * @LastEditTime: 2023-09-24 23:39:48
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

## การแนะนำ
SD-WEBUI-DISCORD นี่คือบอท Discord ที่พัฒนาโดยใช้ Go language สำหรับ [stable-diffusion-webui](https://github.com/AUTOMATIC1111/stable-diffusion-webui). มันใช้ [sd-webui-go](https://github.com/SpenserCai/sd-webui-go) เพื่อเรียกใช้ sd-webui API และรองรับการใช้งานคลัสเตอร์ของโหนด sd-webui หลายโหนดพร้อมการกำหนดเวลาและการจัดสรรอัตโนมัติ
ขณะเดียวกันก็ยังมี [sd-webui-discord-ex](https://github.com/SpenserCai/sd-webui-discord-ex), ซึ่งเป็นส่วนขยายบน stable-diffusion-webui ที่คุณสามารถติดตั้งและใช้งานได้โดยตรง มันจะอัปเดตโดยอัตโนมัติทุกครั้งที่คุณรีสตาร์ท SD WEBUI และรับฟีเจอร์ใหม่ในอนาคต

## ภาพถ่ายหน้าจอ

![First](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/first_page_new.png)

## ข่าว

### **2023-09-24: รองรับการสร้างภาพหลายภาพในคำสั่ง `txt2img` **
### **2023-09-23:**
 - รองรับคำสั่ง `สร้างใหม่` และ `ลบ` ใน `txt2img` 
### **2023-09-22:** 
 - รองรับคำสั่ง `setting_ui`
 - ดีกว่า `txt2img` มี UI ที่ตอบสนอง,ขอขอบคุณ [venetanji](https://github.com/venetanji) ที่สนับสนุนสิ่งนี้! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
 - เพิ่มประสิทธิภาพคำสั่งเพิ่มเมื่อบอทเริ่มทำงาน, ขอขอบคุณ [venetanji](https://github.com/venetanji) ที่สนับสนุนสิ่งนี้! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
### **2023-09-10:** รองรับภาษาจากตำแหน่งที่อยู่ท้องถิ่น
### **2023-09-05:** ศูนย์ผู้ใช้สนับสนุนบน Windows
### **2023-09-04:** รองรับ การสร้างภาพจากภาพ `img2img`
### **2023-08-31:** ศูนย์ผู้ใช้สนับสนุน
### **2023-08-27:**
 - รองรับ txt2img ทางเลือกโมเดล checkpoint
 - รองรับคำสั่ง การอัปโหลดภาพ กับสิ่งที่แนบมาด้วย: `deoldify` `png_info` `roop_image` 
   <details>
   <summary>ดูภาพตัวอย่าง</summary>

      ![example](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/support_attechment.png)
    
   </details>
  
### **2023-08-26: สนับสนุนการตั้งค่าตัวเลือก sd-webui เริ่มต้น**
### **2023-08-23: สนับสนุน ControlNet คำสั่งสำหรับ `txt2img` **
 - โดยใช้ `controlnet_detect` คำสั่งเพื่อรับพารามิเตอร์ของ ControlNet และกรอกลงในไฟล์ `controlnet_args` พารามิเตอร์ของคำสั่ง `txt2img` คุณสามารถใช้ ControlNet ได้ txt2img.
### **2023-08-22: รองรับคำสั่ง `txt2img` **
### **2023-08-22: รองรับคำสั่ง `roop` ** 
### **2023-08-20: รองรับคำสั่ง `controlnet_detect` **  

## คุณสมบัติ
- รองรับภาษาท้องถิ่น
- ตัวเลือกเริ่มต้น sd-webui ทั่วโลก
- รองรับการใช้งานหลายโหนด (sd-webui) คิวคลัสเตอร์แบบกระจายพร้อมการตั้งเวลาอัตโนมัติ
## ศูนย์ผู้ใช้
- ตั้งค่าตัวเลือกผู้ใช้ด้วย UI
- สามารถเปิดใช้งาน User Center ได้อย่างอิสระ - ประเภทฐานข้อมูลที่ปรับแต่งได้ ปัจจุบันรองรับ MySQL และ SQLite
- ผู้ใช้สามารถตั้งค่าตัวเลือกเริ่มต้นของตนเองได้ หากไม่ได้ระบุเมื่อสร้างภาพ ตัวเลือกที่กำหนดค่าของผู้ใช้จะถูกใช้เป็นค่าเริ่มต้น เช่น ขนาดภาพ, cfg, ขั้นตอน ฯลฯ สำหรับรายละเอียดเพิ่มเติม โปรดดูที่ศูนย์ผู้ใช้ (wiki)
- รองรับการลงทะเบียนผู้ใช้
- ดูตัวอย่าง ControlNet
- รองรับการระบุโมดูลและรุ่นผ่านการเลือก ทำให้ไม่จำเป็นต้องป้อนข้อมูลด้วยตนเอง
- อนุญาตให้ดูตัวอย่างเอฟเฟกต์ก่อนการประมวลผลและรับ args พร้อมกัน (ใช้สำหรับที่ผู้ใช้สร้างขึ้น txt2img).
## ข้อความเป็นรูปภาพ
- รองรับการปรับแต่งของ SDXL!
- รองรับตัวเลือกเริ่มต้นที่ผู้ใช้กำหนดจากศูนย์ผู้ใช้
- ช่วยให้สามารถระบุรุ่น ตัวเก็บตัวอย่าง และพารามิเตอร์เสริมอื่นๆ ผ่านการเลือก ซึ่งช่วยลดความจำเป็นในการป้อนข้อมูลด้วยตนเอง
- รองรับการใช้พารามิเตอร์จาก "ControlNet Preview" โดยตรง
- รูปภาพต่อรูปภาพ
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL
- รองรับการทำงานทั้งหมดของ img2img ใน sd-webui!
## Roop Face Swap
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL.
- รองรับการระบุแหล่งที่มาและเป้าหมายผ่านการเลือก ทำให้ไม่จำเป็นต้องป้อนข้อมูลด้วยตนเอง
- รองรับอัลกอริธึมการเรนเดอร์ใบหน้าแบบกำหนดเอง (GFPGAN, CodeFormer) รวมถึงน้ำหนัก
- ขจัดการทำให้สีจางลง
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL
- ปัจจุบันเป็นโมเดลการปรับสีภาพถ่ายที่ดีที่สุด
## SAM-อะไรก็ได้
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL
- รองรับการแบ่งส่วนภาพตามคำอธิบายภาษาธรรมชาติ (DION+SAM)
- ช่วยให้สามารถระบุรุ่น DION และ SAM ผ่านการเลือก ทำให้ไม่จำเป็นต้องป้อนข้อมูลด้วยตนเอง
- การลบพื้นหลัง
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL
- รองรับอัลกอริธึมการลบพื้นหลังที่ใช้กันทั่วไป
- รองรับการคืนหน้ากาก
- เตียงเดี่ยวเสริม
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL
- รองรับการซ่อมแซมใบหน้า
- รองรับความละเอียดสูงสุด โดยมีรุ่นที่มีให้เลือกโดยตรงโดยไม่ต้องป้อนข้อมูลด้วยตนเอง
## ข้อมูลรูปภาพ
- อนุญาตให้อัปโหลดรูปภาพผ่านตัวควบคุมรูปภาพ แทนที่จะใช้ URL 
- รองรับการดึงพารามิเตอร์ของภาพที่สร้างโดย sd-webui

## การใช้งาน 
คำสั่งยังอยู่ระหว่างการพัฒนา และมีสองวิธีในการรับประสบการณ์ sd-webui-discord: 
1. เข้าร่วมเซิร์ฟเวอร์ Discord ของเราซึ่งคุณสามารถลองใช้คุณสมบัติล่าสุดและมีส่วนร่วมโดยการส่งปัญหาและดึงคำขอ
2. ปรับใช้ด้วยตนเองเพื่อให้มีอินสแตนซ์ sd-webui-discord ของคุณเอง

### Discord Server
[![Discord](https://invidget.switchblade.xyz/uNJpzEE4sZ)](https://discord.gg/uNJpzEE4sZ)

<!--

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

</details> -->

### การติดตั้งใช้งาน
คุณต้องติดตั้งส่วนขยายต่อไปนี้บน SD webui:

[sd-webui-segment-anything](https://github.com/continue-revolution/sd-webui-segment-anything)

[sd-weubi-deoldify](https://github.com/SpenserCai/sd-webui-deoldify)

[stable-diffusion-webui-rembg](https://github.com/AUTOMATIC1111/stable-diffusion-webui-rembg)

[sd-webui-roop](https://github.com/s0md3v/sd-webui-roop)

[sd-webui-controlnet](https://github.com/Mikubill/sd-webui-controlnet)

***

1.ดาวน์โหลด เวอร์ชั่นล่าสุดจาก [ที่นี่](https://github.com/SpenserCai/sd-webui-discord/releases/latest).

2.สร้างบัญชีบอทบน Discord และรับโทเค็น [วิธีสร้างแอป Discord](https://discord.com/developers/docs/getting-started).

3.การกำหนดค่าและการเริ่มต้น
```bash
tar -zxvf sd-webui-discord-release-v*.tar.gz # unzip the release package

cd sd-webui-discord-release-v*/release/
```
แก้ไขไฟล์ `config.json` และกรอกโทเค็นและข้อมูลอื่นๆ

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

ถ้าคุณต้องการเซ็ตค่าเริ่มต้น sd-webui
```json
{
    "sd_webui":{
        "servers":[...],
        "default_setting": {
            "cfg_scale": 8,
            "negative_prompt": "bad,text,watermask",
            "height":1024,
            "width":1024,
            "steps":32,
            "sampler":"Euler",
            "sd_model_checkpoint":"sd_xl_base_1.0.safetensors [31e35c80fc]"
        }
    }
    ...
}
```

ถ้าคุณต้องการเปิดใช้งาน **ศูนย์ข้อมูลผู้ใช้**
```json
{
  ...
  "user_center":{
        "enable":false,
        "db_config":{
            "type":"sqlite", // support mysql and sqlite
            "dsn":"./user_center.db"
        }
  }
  ...
}
```

หากคุณต้องการปิดการใช้งานข้อมูลการคืนใน `img2img` and `txt2img`
```json
{
  ...
  "disable_return_gen_info":true
  ...
}
```

เปิดใช้งาน Bot
```bash
# ถ้าคุณไม่สามารถเชื่อมต่อ discord ได้ คุณต้องใช้พรอกซีและรันคำสั่งนี้:
# export https_proxy=http://127.0.0.1:8888;export http_proxy=http://127.0.0.1:8888;
./sd-webui-discord
```

## เข้าร่วม Discord กับเรา
นี่เป็นโครงการที่กำลังดำเนินอยู่ และหากคุณสนใจที่จะมีส่วนร่วม คุณสามารถเข้าร่วมกับเราได้r [Discord Server](https://discord.gg/uNJpzEE4sZ). เรายินดีรับข้อเสนอแนะหรือข้อเสนอแนะ ดังนั้นโปรดส่งปัญหาได้เลย.
