<!--
 * @Author: SpenserCai
 * @Date: 2023-08-17 18:23:21
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 14:22:55
 * @Description: file content
-->
<div align="center">

<img src="https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/logo.png" width="200" height="200" alt="sd-webui-discord">

# SD-WEBUI-DISCORD
Soporte para bots de Discord de Stable Diffusion en clústeres

[LÉEME en múltiples idiomas](https://github.com/SpenserCai/sd-webui-discord/tree/main/docs)
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

## Introducción
SD-WEBUI-DISCORD es un bot de Discord desarrollado en Go para [stable-diffusion-webui](https://github.com/AUTOMATIC1111/stable-diffusion-webui). Utiliza [sd-webui-go](https://github.com/SpenserCai/sd-webui-go) para invocar el API de sd-webui, y soporta el despliegue en clústeres de múltiples nodos de sd-webui con planificación y alojamiento automáticos.
Asimismo, existe también [sd-webui-discord-ex](https://github.com/SpenserCai/sd-webui-discord-ex), que es una extensión de stable-diffusion-webui que puedes instalar y usar directamente. Se actualizará automáticamente cada vez que reinicies la webui de SD.
## Capturas de pantalla

![First](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/first_page_new.png)

## Últimos cambios

### **2023-09-24: Soporte para generación múltiple con el comando `txt2img`**
### **2023-09-23:**
 - Soporte de `Retry` (Reintentar) y `Delete` (Eliminar) en el comando `txt2img`
### **2023-09-22:** 
 - Soporte para el comando `setting_ui`
 - Mejora de la interfaz gráfica de respuesta `txt2img` , gracias a [venetanji](https://github.com/venetanji) por su contribución! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
 - Optimización de añadido de comandos al iniciar el bot, gracias a [venetanji](https://github.com/venetanji) por su contribución! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
### **2023-09-10:** Soporte de idioma local
### **2023-09-05:** Soporte del Centro de Usuarios en Windows
### **2023-09-04:** Soporte de Image to Image
### **2023-08-31:** Soporte del Centro de Usuarios
### **2023-08-27:**
 - Soporte de establecimiento por defecto de cualquier modelo de txt2img
 - Soporte para subir imagen mediante archivo adjunto para los comandos `deoldify` `png_info` `roop_image`
   <details>
   <summary>Ver imagen</summary>

      ![example](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/support_attechment.png)
    
   </details>
  
### **2023-08-26: Soporte para establecer ajustes por defecto de sd-webui**
### **2023-08-23: Soporte de ControlNet para el comando `txt2img`**
 - Puedes usar ControlNet en txt2img mediante el comando `controlnet_detect` para obtener los parámetros de ControlNet, y después, copiando y pegando dichos parámetros en el parámetro `controlnet_args` del comando `txt2img`
### **2023-08-22: Soporte del comando `txt2img`**
### **2023-08-22: Soporte del comando `roop`** 
### **2023-08-20: Soporte del comando `controlnet_detect`**  

## Características
- Soporte de idioma local.
- Ajustes globales por defecto de sd-webui.
- Soporta el desplegamiento multi-nodo de sd-webui, con cola distribuida en clústeres con planificación automática.
- Centro de usuarios
    - Mediante interfaz gráfica, se pueden establecer ajustes personales. 
    - Se puede activar o desactivar el centro del usuarios.
    - Utiliza base de datos personalizable. De momento, cuenta con soporte de MySQL y SQLite.
    - Los usuarios pueden establecer sus propios ajustes personales. Si no se especifican ajustes al generar imágenes, se utilizarán los configurados por defecto, como por ejemplo, las dimensiones de imagen, la cantidad de pasos, o el CFG. Para más detalles, consulte la wiki del Centro de usuarios.
    - Soporta registro de usuarios.
- Vista previa de ControlNet
    - Soporta la selección de módulo y modelo, eliminando así la necesidad de ajuste manual.
    - Permite la vista previa de efectos de preprocesado y obtención de argumentos simultáneamente (los cuales son usados para la generación txt2img por usuario)
- Text to Image
    - Compatible con el refinador de SDXL!
    - Soporta ajustes por defecto preestablecidos por el usuario en el Centro de Usuarios.
    - Permite seleccionar modelo, muestreador, y otros parámetros opcionales, eliminando la necesidad de ajuste manual.
    - Soporta el uso de parámetros desde "Vista previa de ControlNet" directamente.
- Image to Image
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Todas las funciones de img2img de sd-webui son compatibles!
- Intercambio de rostro de Roop
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Soporta la especificación y selección de la imagen fuente y objetivo, eliminando la necesidad de ajuste manual.
    - Soporta algoritmos de renderizado personalizados de rostro (GFPGAN, CodeFormer), así como modelos.
- Deoldify Colorization
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Hasta ahora, el mejor modelo para colorización de fotografías.
- Segment-Anything
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Soporta la segmentación de imágenes basada en descripciones en lenguaje natural (DION+SAM).
    - Soporta la selección de modelos DION y SAM, eliminando la necesidad de ajuste manual.
- Eliminación de fondos
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Compatible con los algoritmos de eliminación de fondos más comúnmente usados.
    - Soporta la devolución de la máscara usada en el proceso.
- Extra Single
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Soporta la reparación facial en imágenes.
    - Soporta superresolución, con modelos disponibles para seleccionar sin tener que hacer ajustes manuales.
- PNG Info
    - Permite subir imágenes mediante archivo adjunto, en lugar de pegar URLs.
    - Soporta la recuperación de parámetros de imágenes que hayan sido generadas por sd-webui.

## Uso

Este bot está todavía en desarrollo, y existen dos maneras para experimentar con sd-webui-discord:
1. Únete a nuestro servidor de Discord donde puedes probar las últimas características y contribuir al proyecto, reportando problemas, aportando sugerencias y haciendo Pull Requests.
2. Despliega el proyecto tú mismo para tener tu propia instancia de sd-webui-discord.

### Servidor de Discord
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

### Instalación
Debes instalar las siguientes extensiones en tu instalación de SD Webui:

[sd-webui-segment-anything](https://github.com/continue-revolution/sd-webui-segment-anything)

[sd-weubi-deoldify](https://github.com/SpenserCai/sd-webui-deoldify)

[stable-diffusion-webui-rembg](https://github.com/AUTOMATIC1111/stable-diffusion-webui-rembg)

[sd-webui-roop](https://github.com/s0md3v/sd-webui-roop)

[sd-webui-controlnet](https://github.com/Mikubill/sd-webui-controlnet)

***

1.Descarga la última versión [aquí](https://github.com/SpenserCai/sd-webui-discord/releases/latest).

2.Crea una nueva cuenta de bot de Discord y obtén su token secreto [Cómo crear una aplicación de Discord](https://discord.com/developers/docs/getting-started).

3.Configuración e inicio
```bash
tar -zxvf sd-webui-discord-release-v*.tar.gz # desempaquetar la última versión

cd sd-webui-discord-release-v*/release/
```
Edita el archivo `config.json` y pega el token secreto y demás información.

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
        "token":"<tu token secreto aquí>",
        "server_id":"<tu id de servidor aquí, si no lo sabes, es opcional>"
    }
}
```

Si quieres establecer ajustes por defecto de sd-webui:
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

Si quieres activar el **Centro de Usuarios**
```json
{
  ...
  "user_center":{
        "enable":false,
        "db_config":{
            "type":"sqlite", //compatibilidad mysql y sqlite
            "dsn":"./user_center.db"
        }
  }
  ...
}
```

Si quieres desactivar la retroalimentación en `img2img` y `txt2img`
```json
{
  ...
  "disable_return_gen_info":true
  ...
}
```

Inicia el bot
```bash
# Si no puedes conectar el bot a discord, debes usar un proxy y ejecutar este comando:
# export https_proxy=http://127.0.0.1:8888;export http_proxy=http://127.0.0.1:8888;
./sd-webui-discord
```

## Contribución
Este es un proyecto en constante desarrollo, y si estás interesado en contribuir, puedes unirte a nuestro [Servidor de Discord](https://discord.gg/uNJpzEE4sZ). Agradecemos toda reseña o sugerencia, asi que siéntete libre de reportar cualquier problema que tengas.

