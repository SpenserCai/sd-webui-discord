/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 00:44:10
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-03 16:38:12
 * @Description: file content
 */
package db_backend

type DbBackend struct{}

// 索引为id字段
type UserInfo struct {
	ID           string `gorm:"column:id;primaryKey;size:100;index:idx_user,unique"`
	Name         string `gorm:"column:name;size:50"`
	Avatar       string `gorm:"column:avatar;size:0"`
	Created      string `gorm:"column:created;size:50"`
	Enable       bool   `gorm:"column:enable"`
	Roles        string `gorm:"column:roles;size:0"`
	StableConfig string `gorm:"column:stable_config;size:0"`
}

type History struct {
	MessageID   string `gorm:"column:message_id;size:100;index:idx_message;index:idx_message_id"`
	UserID      string `gorm:"column:user_id;size:100;index:idx_message"`
	CommandName string `gorm:"column:command_name;size:50"`
	OptionJson  string `gorm:"column:option_json;size:0"`
	Images      string `gorm:"column:images;size:0"`
	Created     string `gorm:"column:created;size:50"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (History) TableName() string {
	return "history"
}
