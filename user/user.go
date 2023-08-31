/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 20:38:24
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 10:40:33
 * @Description: file content
 */
package user

import (
	"encoding/json"

	"reflect"

	"github.com/SpenserCai/sd-webui-discord/user/db"
	"github.com/SpenserCai/sd-webui-discord/user/db/db_backend"
)

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
	Roles        string       `json:"roles"`
	Created      string       `json:"created"`
	StableConfig StableConfig `json:"stable_config"`
}

type UserCenterService struct {
	Db *db.BotDb
}

func (ucs *UserCenterService) GetUserInfo(id string) (*UserInfo, error) {
	userInfo := &db_backend.UserInfo{}
	ucs.Db.Db.Where("id = ?", id).First(userInfo)

	stableConfig := StableConfig{}
	err := json.Unmarshal([]byte(userInfo.StableConfig), &stableConfig)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		Enable:       userInfo.Enable,
		Name:         userInfo.Name,
		Id:           userInfo.ID,
		Roles:        userInfo.Roles,
		StableConfig: stableConfig,
	}, nil
}

func (ucs *UserCenterService) RegisterUser(user *UserInfo) (string, error) {
	// 判断用户是否存在
	userInfo, err := ucs.GetUserInfo(user.Id)
	if err != nil {
		return "", err
	}
	// 如果用户存在则更新用户信息
	if userInfo != nil {
		err := ucs.UpdateUserInfo(user)
		if err != nil {
			return "", err
		}
		return "UPDATED USER INFO", nil
	}

	// 如果用户不存在则创建用户
	newUserInfo := &db_backend.UserInfo{
		ID:           user.Id,
		Name:         user.Name,
		Enable:       user.Enable,
		Roles:        user.Roles,
		StableConfig: "{}",
	}
	err = ucs.Db.Db.Create(newUserInfo).Error
	if err != nil {
		return "", err
	}
	return "REGISTERED USER INFO", nil

}

func (ucs *UserCenterService) UpdateUserInfo(user *UserInfo) error {
	userInfo := &db_backend.UserInfo{
		ID:     user.Id,
		Name:   user.Name,
		Enable: user.Enable,
		Roles:  user.Roles,
	}
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", user.Id).Updates(userInfo).Error
	return err
}

func (ucs *UserCenterService) BanUser(id string) error {
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", id).Update("enable", false).Error
	return err
}

func (ucs *UserCenterService) UpdateStableConfig(user *UserInfo) error {
	return nil
}

func (ucs *UserCenterService) GetUserStableConfigItem(id string, key string, defaultValue interface{}) (interface{}, error) {
	userInfo, err := ucs.GetUserInfo(id)
	if err != nil {
		return nil, err
	}

	stableConfig := userInfo.StableConfig
	structType := reflect.TypeOf(stableConfig)
	structValue := reflect.ValueOf(stableConfig)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")

		if tag == key {
			fieldValue := structValue.Field(i)
			if fieldValue.IsValid() && !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
				return fieldValue.Interface(), nil
			}
			break
		}
	}

	return defaultValue, nil

}
