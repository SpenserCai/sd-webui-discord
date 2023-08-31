/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 00:44:10
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 16:17:28
 * @Description: file content
 */
package db_backend

type DbBackend struct{}

type UserInfo struct {
	ID           string `gorm:"column:id;primaryKey;size:100"`
	Name         string `gorm:"column:name;size:50"`
	Created      string `gorm:"column:created;size:50"`
	Enable       bool   `gorm:"column:enable"`
	Roles        string `gorm:"column:roles;size:0"`
	StableConfig string `gorm:"column:stable_config;size:0"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
