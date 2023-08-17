/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 00:30:18
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 00:44:40
 * @Description: file content
 */
package utils

import (
	"bytes"
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
