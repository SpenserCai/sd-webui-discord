/*
 * @Author: SpenserCai
 * @Date: 2023-08-16 11:05:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-17 13:39:42
 * @Description: file content
 */
package config

type ServerItem struct {
	Name          string `json:"name"`
	Host          string `json:"host"`
	MaxConcurrent int    `json:"max_concurrent"`
	MaxQueue      int    `json:"max_queue"`
	MaxVRAM       string `json:"max_vram"`
}

type Config struct {
	SDWebUi struct {
		Servers []ServerItem `json:"servers"`
	} `json:"sd_webui"`
	Discord struct {
		Token    string `json:"token"`
		ServerId string `json:"server_id"`
	} `json:"discord"`
}
