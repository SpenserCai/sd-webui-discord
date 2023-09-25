/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 00:30:18
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-25 20:57:31
 * @Description: file content
 */
package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"strings"

	"golang.org/x/image/draw"
)

func GetImageBase64(url string) (string, error) {
	if IsUrl(url) {
		imageData, err := GetImageBytesFromUrl(url)
		if err != nil {
			return "", err
		}
		base64Str := ConvertBytesToBase64(imageData)
		return base64Str, nil
	}
	return url, nil
}

func GetImageReaderByBase64(b64 string) (*bytes.Reader, error) {
	imageData, err := ConvertBase64ToBytes(b64)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(imageData), nil
}

func GetImageSizeFromBase64(base64Str string) (int, int, error) {
	// 去除图片标识
	trimmedStr := strings.TrimPrefix(base64Str, "data:image/jpeg;base64,")
	trimmedStr = strings.TrimPrefix(trimmedStr, "data:image/png;base64,")

	// 从Base64字符串解码图片数据
	data, err := base64.StdEncoding.DecodeString(trimmedStr)
	if err != nil {
		return 0, 0, fmt.Errorf("can't decode base64: %s", err)
	}

	// 创建一个Reader以读取图片数据
	reader := strings.NewReader(string(data))

	// 解码图片
	img, _, err := image.Decode(reader)
	if err != nil {
		return 0, 0, fmt.Errorf("can't decode image: %s", err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	return width, height, nil
}

func ResizeImage(inputImage image.Image, scalePercent float64) (image.Image, error) {
	// 获取原始图像的尺寸
	originalBounds := inputImage.Bounds()
	originalWidth := originalBounds.Dx()
	originalHeight := originalBounds.Dy()

	// 计算缩放后的尺寸
	newWidth := int(float64(originalWidth) * scalePercent / 100)
	newHeight := int(float64(originalHeight) * scalePercent / 100)

	// 创建缩放后的图像
	resizedImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// 使用双线性插值算法进行图像缩放
	draw.CatmullRom.Scale(resizedImage, resizedImage.Bounds(), inputImage, originalBounds, draw.Over, nil)

	// 将缩放后的图像编码为 JPEG 格式
	var outputBuffer bytes.Buffer
	err := jpeg.Encode(&outputBuffer, resizedImage, nil)
	if err != nil {
		return nil, err
	}

	// 解码 JPEG 图像数据，生成最终的缩放后图像
	finalImage, err := jpeg.Decode(&outputBuffer)
	if err != nil {
		return nil, err
	}

	return finalImage, nil
}

// 图片合并
func MergeImageFromBase64(base64List []string) (string, error) {
	// 如果是单张图片，直接返回
	if len(base64List) == 1 {
		return base64List[0], nil
	}
	var completeImage image.Image
	var resizeErr error
	// 如果是两张图片，水平合并
	if len(base64List) == 2 {
		for _, base64Str := range base64List {
			// 去除图片标识
			trimmedStr := strings.TrimPrefix(base64Str, "data:image/jpeg;base64,")
			trimmedStr = strings.TrimPrefix(trimmedStr, "data:image/png;base64,")

			// 从Base64字符串解码图片数据
			data, err := base64.StdEncoding.DecodeString(trimmedStr)
			if err != nil {
				return "", fmt.Errorf("can't decode base64: %s", err)
			}

			// 创建一个Reader以读取图片数据
			reader := strings.NewReader(string(data))

			// 解码图片
			img, _, err := image.Decode(reader)
			if err != nil {
				return "", fmt.Errorf("can't decode image: %s", err)
			}
			if completeImage == nil {
				completeImage = img
			} else {
				// 把图片水平拼接到completeImage上，直接实现算法
				newWidth := completeImage.Bounds().Dx() + img.Bounds().Dx()
				newHeight := completeImage.Bounds().Dy()
				newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
				draw.Draw(newImage, image.Rect(0, 0, completeImage.Bounds().Dx(), newHeight), completeImage, image.Point{0, 0}, draw.Src)
				draw.Draw(newImage, image.Rect(completeImage.Bounds().Dx(), 0, newWidth, newHeight), img, image.Point{0, 0}, draw.Src)
				completeImage = newImage
			}
		}
		// 缩放到原来的75%
		completeImage, resizeErr = ResizeImage(completeImage, 75)
	} else if len(base64List) == 4 { // 如果是4张图片，2x2合并
		var upPartImage image.Image
		var downPartImage image.Image
		for i, base64Str := range base64List {
			// 去除图片标识
			trimmedStr := strings.TrimPrefix(base64Str, "data:image/jpeg;base64,")
			trimmedStr = strings.TrimPrefix(trimmedStr, "data:image/png;base64,")
			// 从Base64字符串解码图片数据
			data, err := base64.StdEncoding.DecodeString(trimmedStr)
			if err != nil {
				return "", fmt.Errorf("can't decode base64: %s", err)
			}
			reader := strings.NewReader(string(data))
			// 解码图片
			img, _, err := image.Decode(reader)
			if err != nil {
				return "", fmt.Errorf("can't decode image: %s", err)
			}
			// 把第一第二张水平合并，第三第四张水平合并，再把两张图片垂直合并
			if i == 0 || i == 1 {
				if upPartImage == nil {
					upPartImage = img
				} else {
					newWidth := upPartImage.Bounds().Dx() + img.Bounds().Dx()
					newHeight := upPartImage.Bounds().Dy()
					newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
					draw.Draw(newImage, image.Rect(0, 0, upPartImage.Bounds().Dx(), newHeight), upPartImage, image.Point{0, 0}, draw.Src)
					draw.Draw(newImage, image.Rect(upPartImage.Bounds().Dx(), 0, newWidth, newHeight), img, image.Point{0, 0}, draw.Src)
					upPartImage = newImage
				}
			}
			if i == 2 || i == 3 {
				if downPartImage == nil {
					downPartImage = img
				} else {
					newWidth := downPartImage.Bounds().Dx() + img.Bounds().Dx()
					newHeight := downPartImage.Bounds().Dy()
					newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
					draw.Draw(newImage, image.Rect(0, 0, downPartImage.Bounds().Dx(), newHeight), downPartImage, image.Point{0, 0}, draw.Src)
					draw.Draw(newImage, image.Rect(downPartImage.Bounds().Dx(), 0, newWidth, newHeight), img, image.Point{0, 0}, draw.Src)
					downPartImage = newImage
				}
			}
		}
		// 把两张图片垂直合并
		newWidth := upPartImage.Bounds().Dx()
		newHeight := upPartImage.Bounds().Dy() + downPartImage.Bounds().Dy()
		newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		draw.Draw(newImage, image.Rect(0, 0, newWidth, upPartImage.Bounds().Dy()), upPartImage, image.Point{0, 0}, draw.Src)
		draw.Draw(newImage, image.Rect(0, upPartImage.Bounds().Dy(), newWidth, newHeight), downPartImage, image.Point{0, 0}, draw.Src)
		completeImage = newImage
		// 缩放到原来的50%
		completeImage, resizeErr = ResizeImage(completeImage, 50)
	} else { // 其他情况，三张图一行，多行合并，空白部分用白色填充
		var upPartImage image.Image
		var downPartImage image.Image
		var completeImageList []image.Image
		for i, base64Str := range base64List {
			// 去除图片标识
			trimmedStr := strings.TrimPrefix(base64Str, "data:image/jpeg;base64,")
			trimmedStr = strings.TrimPrefix(trimmedStr, "data:image/png;base64,")
			// 从Base64字符串解码图片数据
			data, err := base64.StdEncoding.DecodeString(trimmedStr)
			if err != nil {
				return "", fmt.Errorf("can't decode base64: %s", err)
			}
			reader := strings.NewReader(string(data))
			// 解码图片
			img, _, err := image.Decode(reader)
			if err != nil {
				return "", fmt.Errorf("can't decode image: %s", err)
			}
			// 把图片水平拼接到completeImage上，直接实现算法
			if upPartImage == nil {
				upPartImage = img
			} else {
				newWidth := upPartImage.Bounds().Dx() + img.Bounds().Dx()
				newHeight := upPartImage.Bounds().Dy()
				newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
				draw.Draw(newImage, image.Rect(0, 0, upPartImage.Bounds().Dx(), newHeight), upPartImage, image.Point{0, 0}, draw.Src)
				draw.Draw(newImage, image.Rect(upPartImage.Bounds().Dx(), 0, newWidth, newHeight), img, image.Point{0, 0}, draw.Src)
				upPartImage = newImage
			}
			if i%3 == 2 {
				completeImageList = append(completeImageList, upPartImage)
				upPartImage = nil
			}
		}
		if upPartImage != nil {
			completeImageList = append(completeImageList, upPartImage)
		}
		// 把图片垂直拼接到completeImage上，直接实现算法
		for i, img := range completeImageList {
			if downPartImage == nil {
				downPartImage = img
			} else {
				newWidth := downPartImage.Bounds().Dx()
				newHeight := downPartImage.Bounds().Dy() + img.Bounds().Dy()
				newImage := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
				draw.Draw(newImage, image.Rect(0, 0, newWidth, downPartImage.Bounds().Dy()), downPartImage, image.Point{0, 0}, draw.Src)
				draw.Draw(newImage, image.Rect(0, downPartImage.Bounds().Dy(), newWidth, newHeight), img, image.Point{0, 0}, draw.Src)
				downPartImage = newImage
			}
			if i%3 == 2 {
				completeImage = downPartImage
				downPartImage = nil
			}
		}
		if downPartImage != nil {
			completeImage = downPartImage
		}
		// 缩放到原来的50%
		completeImage, resizeErr = ResizeImage(completeImage, 50)
	}
	if resizeErr != nil {
		return "", resizeErr
	}
	// completeImage 转 base64
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, completeImage, nil)
	if err != nil {
		return "", err
	}
	return ConvertBytesToBase64(buf.Bytes()), nil
}
