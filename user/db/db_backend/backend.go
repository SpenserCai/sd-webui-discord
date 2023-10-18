/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 00:44:10
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-18 21:43:21
 * @Description: file content
 */
package db_backend

type DbBackend struct{}

// 索引为id字段
type UserInfo struct {
	ID                 string `gorm:"column:id;primaryKey;size:100;index:idx_user,unique"`
	Name               string `gorm:"column:name;size:50"`
	Avatar             string `gorm:"column:avatar;size:0"`
	Created            string `gorm:"column:created;size:50"`
	Enable             bool   `gorm:"column:enable"`
	Roles              string `gorm:"column:roles;size:0"`
	StableConfig       string `gorm:"column:stable_config;size:0"`
	CycleCredit        int32  `gorm:"column:cycle_credit;default:0"`      // 周期credit
	CreditUpdateCycle  string `gorm:"column:credit_update_cycle;size:50"` // credit更新周期 格式 1|H 1|D 1|W 1|M
	CycleCreditUpdated string `gorm:"column:credit_updated;size:50"`      // 周期credit更新时间，用于判断是否需要更新credit
	PlusCredit         int32  `gorm:"column:plus_credit;default:0"`       // 充值的credit
	IsPrivate          bool   `gorm:"column:is_private;default:false"`
}

type History struct {
	MessageID      string `gorm:"column:message_id;size:100;index:idx_message;index:idx_message_id"`
	UserID         string `gorm:"column:user_id;size:100;index:idx_message"`
	CommandName    string `gorm:"column:command_name;size:50"`
	OptionJson     string `gorm:"column:option_json;size:0"`
	Images         string `gorm:"column:images;size:0"`
	Created        string `gorm:"column:created;size:50"`
	Deleted        bool   `gorm:"column:deleted;default:false"`
	ImageBlurHashs string `gorm:"column:image_blur_hashs;size:0"`
	IsPrivate      bool   `gorm:"column:is_private;default:false"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

func (History) TableName() string {
	return "history"
}
