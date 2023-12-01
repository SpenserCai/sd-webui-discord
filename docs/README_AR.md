<div align="center">

<img src="https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/logo.png" width="200" height="200" alt="sd-webui-discord">

# SD-WEBUI-DISCORD
دعم لبوت ديسكورد لواجهة مستخدم ستيبل ديفيوجن

[README بعدة لغات](https://github.com/SpenserCai/sd-webui-discord/tree/main/docs)
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

## المقدمة
SD-WEBUI-DISCORD هو بوت ديسكورد تم تطويره بلغة Go لصالح [واجهة مستخدم مستقرة للتفاعل عبر الويب](https://github.com/AUTOMATIC1111/stable-diffusion-webui). يستفيد من [sd-webui-go](https://github.com/SpenserCai/sd-webui-go) لاستدعاء واجهة برمجة التطبيقات لـ sd-webui ويدعم نشر أكثر من وحدة sd-webui بتكوين تلقائي وتوزيع ذكي.
في الوقت نفسه، يوجد أيضًا [sd-webui-discord-ex](https://github.com/SpenserCai/sd-webui-discord-ex)، وهو امتداد لواجهة المستخدم الثابتة للتفاعل عبر الويب يمكنك تثبيته واستخدامه مباشرة. سيتم تحديثه تلقائيًا في كل مرة تقوم فيها بإعادة تشغيل SD webui.

## لقطات الشاشة

![الأولى](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/first_page_new.png)
![الثانية](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/second_page_new.png)

## الأخبار
 - `12 أكتوبر 2023`: دعم أمر `Website`
   <details>
     <summary>انظر الصورة</summary>
       
      ![مثال](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/website_gallery.png)
      ![مثال](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/website_community.png)
      ![مثال](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/website_image_detail.png)

   </details>
 - `2023-09-24`: دعم إنشاء صور متعددة في أمر `txt2img`
 - `2023-09-23`: دعم `إعادة المحاولة` و `الحذف` في أمر `txt2img`
 - `2023-09-22`:
     - دعم أمر `setting_ui`
     - تحسين واجهة الرد في `txt2img`، شكرًا لـ [venetanji](https://github.com/venetanji) لدعمه! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
     - تحسين إضافة الأوامر عند بدء تشغيل البوت، شكرًا لـ [venetanji](https://github.com/venetanji) لدعمه! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
 - `2023-09-10`: دعم اللغة المحلية
 - `2023-09-05`: دعم مركز المستخدم على نظام Windows
 - `2023-09-04`: دعم تحويل الصورة إلى صورة
 - `2023-08-31`: دعم مركز المستخدم
 - `2023-08-27`:
     - دعم اختيار نقطة فحص النموذج في `txt2img`
     - دعم رفع الصور بواسطة المرفقات: أوامر `deoldify`، `png_info`، `roop_image`
     <details>
     <summary>راجع الصورة</summary>

      ![مثال](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/support_attechment.png)
     
     </details>

  
 - `24 سبتمبر 2023`: دعم إنشاء عدة صور في أمر `txt2img`
 - `23 سبتمبر 2023`: دعم `إعادة المحاولة` و `الحذف` في أمر `txt2img`
 - `22 سبتمبر 2023` 
     - دعم أمر `setting_ui`
     - تحسين واجهة الاستجابة لأمر `txt2img`، شكرًا لـ [venetanji](https://github.com/venetanji) للدعم! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
     - تحسين إضافة الأوامر عند بدء تشغيل البوت، شكرًا لـ [venetanji](https://github.com/venetanji) للدعم! [#5](https://github.com/SpenserCai/sd-webui-discord/pull/5)
 - `10 سبتمبر 2023`: دعم اللغة المحلية
 - `5 سبتمبر 2023`: دعم مركز المستخدم على نظام Windows
 - `4 سبتمبر 2023`: دعم من الصورة إلى الصورة
 - `31 أغسطس 2023`: دعم مركز المستخدم
 - `27 أغسطس 2023`:
     - دعم نقطة فحص نموذج `txt2img`
     - دعم رفع الصور بمرفقات: أوامر `deoldify`، `png_info`، `roop_image`
     <details>
     <summary>انظر الصورة</summary>

      ![مثال](https://raw.githubusercontent.com/SpenserCai/sd-webui-discord/main/res/support_attechment.png)
     
     </details>
  
 - `26 أغسطس 2023`: دعم إعدادات sd-webui الافتراضية
 - `23 أغسطس 2023`: دعم ControlNet لأمر `txt2img`
     - باستخدام أمر `controlnet_detect` للحصول على معلمات ControlNet وملءها في معلمة `controlnet_args` لأمر `txt2img`، يمكنك استخدام ControlNet في `txt2img`.
 - `22 أغسطس 2023`: دعم أمر `txt2img`
 - `22 أغسطس 2023`: دعم أمر `roop`
 - `20 أغسطس 2023`: دعم أمر `controlnet_detect`  

## الميزات
- الموقع الإلكتروني
    - معرض المستخدمين.
    - معرض المجتمع.
    - عرض إعدادات المستخدم.
    - قائمة المستخدمين (للمسؤولين).
    - قائمة عقدة العنقود (للمسؤولين).
- دعم اللغة المحلية.
- خيارات sd-webui الافتراضية على مستوى العالم.
- يدعم نشر متعدد العقد (sd-webui)، طابور عقد موزع مع جدولة تلقائية.
- مركز المستخدم
    - تحديد خيارات المستخدم مع واجهة المستخدم. 
    - يمكن تمكين مركز المستخدم بحرية.
    - أنواع قاعدة البيانات قابلة للتخصيص، حاليًا تدعم MySQL و SQLite.
    - يمكن للمستخدمين تعيين خياراتهم الافتراضية. إذا لم يتم تحديدها عند إنشاء الصور، سيتم استخدام خيارات المستخدم المكونة تلقائيًا، مثل أبعاد الصورة، والتكوين، والخطوات، وما إلى ذلك. لمزيد من التفاصيل، يرجى الرجوع إلى مركز المستخدم (ويكي).
    - يدعم تسجيل المستخدم.
- معاينة ControlNet
    - يدعم تحديد الوحدة والنموذج من خلال التحديد، مما يلغي الحاجة إلى إدخال يدوي.
    - يسمح بمعاينة تأثيرات المعالجة المسبقة والحصول على وسائط في وقت واحد (يستخدم لإنشاء صور txt2img من قبل المستخدم).
- نص إلى صورة
    - دعم مصقل SDXL!
    - يدعم الخيارات الافتراضية المحددة من مركز المستخدم.
    - يمكن تحديد الموديل وأخذ العينات ومعلمات أخرى اختياريّة من خلال التحديد، مما يلغي الحاجة إلى إدخال يدوي. 
    - يدعم استخدام المعلمات من "معاينة ControlNet" مباشرة. 
- صورة إلى صورة
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - يدعم جميع عمليات img2img في sd-webui!
- Roop Face Swap
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - يدعم تحديد المصدر والهدف من خلال التحديد، مما يلغي الحاجة إلى إدخال يدوي.
    - يدعم خوارزميات تجديد الوجه المخصصة (GFPGAN، CodeFormer)، بالإضافة إلى الأوزان.
- تلوين Deoldify
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - حاليًا أفضل نموذج لتلوين الصور الفوتوغرافية.
- Segment-Anything
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - يدعم تقسيم الصورة استنادًا إلى وصف اللغة الطبيعية (DION+SAM).
    - يسمح بتحديد نماذج DION و SAM من خلال التحديد، مما يلغي الحاجة إلى إدخال يدوي.
- إزالة الخلفية
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - يدعم خوارزميات إزالة الخلفية الشائعة.
    - يدعم إرجاع القناع.
- Extra Single
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - يدعم إصلاح الوجه.
    - يدعم التكبير، مع توفر نماذج للاختيار المباشر دون إدخال يدوي.
- Png Info
    - يسمح بتحميل الصور عبر وحدة التحكم في الصور بدلاً من استخدام عناوين URL.
    - يدعم استرداد معلمات الصور التي تم إنشاءها بواسطة sd-webui.

## الاستخدام

الأمر لا يزال قيد التطوير النشط، وهناك طريقتان لتجربة sd-webui-discord:
1. انضم إلى خادم Discord الخاص بنا حيث يمكنك تجربة أحدث الميزات والمساهمة من خلال تقديم مشاكل وطلبات سحب.
2. نصبه بنفسك لتكون لديك نسخة خاصة من sd-webui-discord.

[لوحة معلومات SD-WEBUI-DISCORD](https://aigc.ngrok.io/)

[![Discord](https://invidget.switchblade.xyz/uNJpzEE4sZ)](https://discord.gg/uNJpzEE4sZ)


## الوثائق
للإحالة إلى الدورة الدراسية المفصلة [الويكي](https://github.com/SpenserCai/sd-webui-discord/wiki)


## المشاركة
هذا مشروع مستمر، وإذا كنت مهتمًا بالمساهمة، يمكنك الانضمام إلى [خادم Discord الخاص بنا](https://discord.gg/uNJpzEE4sZ). نحن نرحب بأي تعليق أو اقتراح، لذا لا تتردد في تقديم مشكلة.

## مستخدم بواسطة

### WAN Show Bingo!

[![Discord](https://invidget.switchblade.xyz/pWS5mw7jFz)](https://discord.gg/pWS5mw7jFz)

### AIGC
[![Discord](https://invidget.switchblade.xyz/aigc)](https://discord.gg/aigc)

### MultiPlayer DAO
[![Discord](https://invidget.switchblade.xyz/XsJgWfDqjR)](https://discord.gg/XsJgWfDqjR)

##  لدعم تطوير
إذا أعجبك هذا المشروع وترغب في أن تصبح مؤيدًا أو راعيًا، ستحصل على:

دعم فني فردي.
سيتم عرض شعارك أو شعار شركتك في صفحة الراعي لهذه الصفحة.

<a href="https://www.patreon.com/sd_webui_discord"><img alt="Sponsor with Patreon" title="Sponsor with Patreon" src="https://img.shields.io/badge/-Sponsor-ea4aaa?style=for-the-badge&logo=patreon&logoColor=white"/></a>
