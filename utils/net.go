/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 00:38:37
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 00:40:30
 * @Description: file content
 */
package utils

import (
	"io"
	"net/http"
	"strings"
)

func IsUrl(str string) bool {
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

func GetImageBytesFromUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	imageData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return imageData, nil
}
