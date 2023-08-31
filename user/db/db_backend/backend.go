/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 00:44:10
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 10:11:42
 * @Description: file content
 */
package db_backend

type DbBackend struct{}

type UserInfo struct {
	ID           string `gorm:"column:id;primaryKey"`
	Name         string `gorm:"column:name"`
	Enable       bool   `gorm:"column:enable"`
	Roles        string `gorm:"column:roles;size:0"`
	StableConfig string `gorm:"column:stable_config;size:0"`
}

func (UserInfo) TableName() string {
	return "user_info"
}
