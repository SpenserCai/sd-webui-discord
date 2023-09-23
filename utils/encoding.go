/*
 * @Author: SpenserCai
 * @Date: 2023-08-17 00:40:47
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-24 00:43:58
 * @Description: file content
 */
package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"reflect"
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

// 判断字符串是否是json字符串
func IsJsonString(str string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(str), &js) == nil
}

func IsZeroValue(value interface{}) bool {
	if value == nil {
		return true
	}

	reflectValue := reflect.ValueOf(value)
	zeroValue := reflect.Zero(reflectValue.Type()).Interface()

	return reflect.DeepEqual(value, zeroValue)
}

func GetJsonReaderByJsonString(jsonStr string) (*strings.Reader, error) {
	// 格式化json字符串
	var fJsonStr bytes.Buffer
	err := json.Indent(&fJsonStr, []byte(jsonStr), "", "  ")
	if err != nil {
		return nil, err
	}
	// 把格式化后的json字符串转换成Reader
	jsonReader := strings.NewReader(fJsonStr.String())
	return jsonReader, nil
}
