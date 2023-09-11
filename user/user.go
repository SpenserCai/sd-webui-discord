/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 20:38:24
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-11 14:59:58
 * @Description: file content
 */
package user

import (
	"encoding/json"
	"strings"
	"time"

	"reflect"

	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/user/db"
	"github.com/SpenserCai/sd-webui-discord/user/db/db_backend"
)

// 权限表,如果命令出现再某个role中，则代表只有这个role的用户才能使用这个命令，key是role，value是命令
var PermissionTable = map[string][]string{
	"admin": {"ban", "unban"},
}

type StableConfig struct {
	Model          string  `json:"sd_model_checkpoint"`
	Height         int64   `json:"height"`
	Width          int64   `json:"width"`
	Steps          int64   `json:"steps"`
	CfgScale       float64 `json:"cfg_scale"`
	NegativePrompt string  `json:"negative_prompt"`
	Sampler        string  `json:"sampler"`
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

func NewUserCenterService(ucsCfg *config.UserCenter) (*UserCenterService, error) {
	db, err := db.NewBotDb(&ucsCfg.DbConfig)
	if err != nil {
		return nil, err
	}

	return &UserCenterService{
		Db: db,
	}, nil
}

func (ucs *UserCenterService) GetUserInfo(id string) (*UserInfo, error) {
	userInfo := &db_backend.UserInfo{}
	ucs.Db.Db.Where("id = ?", id).First(userInfo)
	if userInfo.ID == "" {
		return nil, nil
	}

	stableConfig := StableConfig{}
	if userInfo.StableConfig != "" {
		err := json.NewDecoder(strings.NewReader(userInfo.StableConfig)).Decode(&stableConfig)
		if err != nil {
			return nil, err
		}
	}

	return &UserInfo{
		Enable:       userInfo.Enable,
		Name:         userInfo.Name,
		Id:           userInfo.ID,
		Roles:        userInfo.Roles,
		StableConfig: stableConfig,
	}, nil
}

func (ucs *UserCenterService) CheckUserPermission(id string, cmd string) bool {
	userInfo, err := ucs.GetUserInfo(id)
	if err != nil {
		return false
	}
	if userInfo == nil {
		// 判断命令是否再任意一个role中，如果在返回false
		for _, permissionRoles := range PermissionTable {
			for _, permissionRole := range permissionRoles {
				if permissionRole == cmd {
					return false
				}
			}
		}
		return true
	}
	roles := strings.Split(userInfo.Roles, ",")
	returnValue := true
	// 判断命令是否再任意一个role中，如果在判断用户是否有这个role，如果有返回true，如果没有返回false
	for cRole, cmds := range PermissionTable {
		for _, cCmd := range cmds {
			if cCmd == cmd {
				for _, role := range roles {
					if role == cRole {
						return true
					}
				}
				returnValue = false
			}
		}
	}
	return returnValue
}

func (ucs *UserCenterService) RegisterUser(user *UserInfo) (string, error) {
	// 判断用户是否存在
	userInfo, err := ucs.GetUserInfo(user.Id)
	if err != nil {
		return "", err
	}
	// 如果用户存在则更新用户信息
	if userInfo != nil {
		if !userInfo.Enable {
			return "USER BANNED", nil
		}
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
		Created:      time.Now().Format("2006-01-02 15:04:05"),
		Enable:       true,
		Roles:        "user",
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

	stableConfig, err := json.Marshal(user.StableConfig)
	if err != nil {
		return err
	}
	err = ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", user.Id).Update("stable_config", stableConfig).Error
	return err
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
