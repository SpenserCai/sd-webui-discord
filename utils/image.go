/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 00:30:18
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-04 14:35:21
 * @Description: file content
 */
package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"strings"
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
