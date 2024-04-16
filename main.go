/*
 * @Author: SpenserCai
 * @Date: 2023-08-15 21:55:36
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-08 17:00:50
 * @Description: file content
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/SpenserCai/sd-webui-discord/api"
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

func RunWebSite() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(exePath)
	websiteDir := filepath.Join(exeDir, "website")
	fs := http.FileServer(http.Dir(websiteDir))
	http.Handle("/", fs)
	log.Println("website dir:", websiteDir)
	apiURL, err := url.Parse(fmt.Sprintf("http://%s:%d", global.Config.WebSite.Api.Host, global.Config.WebSite.Api.Port))
	if err != nil {
		return err
	}
	http.Handle("/api/", httputil.NewSingleHostReverseProxy(apiURL))
	log.Println("api url:", apiURL)
	// 启动Web服务器
	addr := fmt.Sprintf("%s:%d", global.Config.WebSite.Web.Host, global.Config.WebSite.Web.Port) // 替换为实际的端口号
	log.Printf("WebSite started on port %s\n", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}
	return nil
}

func OpenWebSite(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "windows": // Windows
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // Linux 或其他 Unix 系统
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
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
	global.Features = global.InitFeatureGate()
	disBot, err := dbot.NewDiscordBot(global.Config.Discord.Token, global.Config.Discord.ServerId)
	if err != nil {
		log.Println(err)
		return
	}
	go api.StartApiService()
	if global.Config.WebSite.Web.StartWithServer {
		go RunWebSite()
		if global.Config.WebSite.Web.OpenBrowser {
			// 使用默认浏览器打开网页
			err = OpenWebSite(fmt.Sprintf("http://%s:%d", global.Config.WebSite.Web.Host, global.Config.WebSite.Web.Port))
			if err != nil {
				log.Println(err)
			}
		}

	}
	disBot.Run()

}
