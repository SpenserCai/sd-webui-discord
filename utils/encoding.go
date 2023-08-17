/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 00:40:47
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 11:19:55
 * @Description: file content
 */
package utils

import (
	"encoding/base64"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ConvertBytesToBase64(data []byte) string {
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str
}

func ConvertBase64ToBytes(base64Str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(base64Str)
	return data, err
}

// 如果name中有_则用下划线分割后每个首字母专大写，如果没有_则直接首字母转大写
func FormatCommand(cmdName string) string {
	cmdSplit := strings.Split(cmdName, "_")
	reCmdName := ""
	for _, v := range cmdSplit {
		reCmdName += cases.Title(language.English, cases.NoLower).String(v)
	}
	return reCmdName
}
