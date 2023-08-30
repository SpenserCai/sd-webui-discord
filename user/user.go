/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 20:38:24
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-30 21:19:51
 * @Description: file content
 */
package user

type StableConfig struct {
	Model          string  `json:"sd_model_checkpoint"`
	Height         int64   `json:"height"`
	Width          int64   `json:"width"`
	Steps          int64   `json:"steps"`
	CfgScale       float64 `json:"cfg_scale"`
	NegativePrompt string  `json:"negative_prompt"`
}

type UserInfo struct {
	Enable       bool         `json:"enable"`
	Name         string       `json:"name"`
	Id           string       `json:"id"`
	StableConfig StableConfig `json:"config"`
}

type UserCenterService struct{}

func (ucs *UserCenterService) GetUserInfo(id string) *UserInfo {
	return &UserInfo{}
}

func (ucs *UserCenterService) RegisterUser(user *UserInfo) error {
	return nil
}

func (ucs *UserCenterService) BanUser(user *UserInfo) error {
	return nil
}

func (ucs *UserCenterService) UpdateStableConfig(user *UserInfo) error {
	return nil
}

func (ucs *UserCenterService) GetUserStableConfigItem(user *UserInfo, key string) interface{} {
	return nil
}
