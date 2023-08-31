/*
 * @Author: SpenserCai
 * @Date: 2023-08-15 21:55:36
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 11:48:39
 * @Description: file content
 */
package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/SpenserCai/sd-webui-discord/cluster"
	"github.com/SpenserCai/sd-webui-discord/dbot"
	"github.com/SpenserCai/sd-webui-discord/global"
	"github.com/SpenserCai/sd-webui-discord/queue"
	"github.com/SpenserCai/sd-webui-discord/user"
)

func LoadConfig() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, "config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&global.Config)
	if err != nil {
		return err
	}
	return nil
}

func InitClusterManager() {
	global.ClusterManager = cluster.NewClusterService(global.Config)
	global.ClusterManager.Start()
}

func InitUserCenterService() error {
	var err error
	global.UserCenterSvc, err = user.NewUserCenterService(&global.Config.UserCenter)
	return err

}

func PrintEvent() {
	for {
		event := global.ClusterManager.GetEvent()
		eventName := "Event"
		switch event.EventType {
		case queue.EventPendding:
			eventName = "Pendding"
		case queue.EventRunning:
			eventName = "Running"
		case queue.EventSuccess:
			eventName = "Success"
		case queue.EventFaile:
			eventName = "Failed"
		case queue.EventCancel:
			eventName = "Cancel"
		default:
			eventName = "Unknown"
		}

		log.Printf("[Event]: ID:%v Type:%v", event.ID, eventName)
	}
}

func main() {
	err := LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	if global.Config.UserCenter.Enable {
		err := InitUserCenterService()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("UserCenterService Init Success")
	}
	InitClusterManager()
	go PrintEvent()
	disBot, err := dbot.NewDiscordBot(global.Config.Discord.Token, global.Config.Discord.ServerId)
	if err != nil {
		log.Println(err)
		return
	}
	disBot.Run()

}
