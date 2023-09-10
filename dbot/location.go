/*
 * @Author: SpenserCai
 * @Date: 2023-09-10 16:18:27
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-10 22:17:50
 * @Description: file content
 */
package dbot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/bwmarrin/discordgo"
)

type LocationItem struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type Location struct {
	LocationList []LocationItem `json:"location_list"`
}

func GetLocationItem(filePath string) (string, []LocationItem, error) {
	var location Location
	var locationItem []LocationItem
	// 从filePath中读取文件名，并将文件 内容读取到locationItem中 不用扩展名
	locationName := filepath.Base(filePath)
	locationName = locationName[:len(locationName)-len(filepath.Ext(locationName))]
	file, err := os.Open(filePath)
	if err != nil {
		return "", locationItem, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&location)
	if err != nil {
		return "", locationItem, err
	}
	return locationName, location.LocationList, nil

}

func GetLocation(locationName string) discordgo.Locale {
	// 循环discordgo.Locales找到Value为locationName的Locale
	keys := reflect.ValueOf(discordgo.Locales).MapKeys()

	for _, key := range keys {
		locale := key.Interface().(discordgo.Locale)
		if string(locale) == locationName {
			return locale
		}
	}
	return discordgo.Unknown

}

func AddLocationDescriptionMap(cmd *discordgo.ApplicationCommand) {
	// 如果cmd的DescriptionLocalizations是零值，则初始化
	if cmd.DescriptionLocalizations == nil {
		cmd.DescriptionLocalizations = &map[discordgo.Locale]string{}
	}

}

func (dbot *DiscordBot) SetLocation() error {
	locationMap := make(map[string][]LocationItem)
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath) + "/location"
	// 将exeDir中的所有文件名读取到locationList中
	err = filepath.Walk(exeDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			name, locationItem, err := GetLocationItem(filepath.Join(exeDir, info.Name()))
			if err != nil {
				return err
			}
			locationMap[name] = locationItem
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	// log.Println(locationMap)
	// 循环遍历dbot.AppCommand
	for _, cmd := range dbot.AppCommand {
		AddLocationDescriptionMap(cmd)
		// 循环遍历locationMap
		for k, localCmdItemList := range locationMap {
			locale := GetLocation(k)
			if locale == discordgo.Unknown {
				continue
			}
			for _, localCmdItem := range localCmdItemList {
				if cmd.Name == localCmdItem.Command && localCmdItem.Description != "" {
					(*cmd.DescriptionLocalizations)[locale] = localCmdItem.Description
				}
			}

		}
	}
	return nil
}
