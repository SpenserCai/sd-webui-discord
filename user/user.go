/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 20:38:24
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-09-29 12:03:00
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
	"admin": {"user_info", "cluster_status"},
}

type StableConfig struct {
	Model          string  `json:"sd_model_checkpoint"`
	Vae            string  `json:"sd_vae"`
	Height         int64   `json:"height"`
	Width          int64   `json:"width"`
	Steps          int64   `json:"steps"`
	CfgScale       float64 `json:"cfg_scale"`
	NegativePrompt string  `json:"negative_prompt"`
	Sampler        string  `json:"sampler"`
	ClipSkip       int64   `json:"CLIP_stop_at_last_layers"`
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
		ID:      user.Id,
		Name:    user.Name,
		Created: time.Now().Format("2006-01-02 15:04:05"),
		Enable:  true,
		// 如果用户总数为0，则创建的用户为admin，否则为user
		Roles: func() string {
			count, err := ucs.GetUserCount()
			if err == nil && count == 0 {
				return "user,admin"
			}
			return "user"
		}(),
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

func (ucs *UserCenterService) UnBanUser(id string) error {
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", id).Update("enable", true).Error
	return err
}

func (ucs *UserCenterService) IsBaned(id string) (bool, error) {
	userInfo, err := ucs.GetUserInfo(id)
	if err != nil {
		return false, err
	}
	if userInfo == nil {
		return false, nil
	}
	return !userInfo.Enable, nil
}

func (ucs *UserCenterService) IsRegistered(id string) bool {
	userInfo, err := ucs.GetUserInfo(id)
	if err != nil {
		return false
	}
	if userInfo == nil {
		return false
	}
	return true
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

// 写入用户历史记录
func (ucs *UserCenterService) WriteUserHistory(messageId string, userId string, commandName string, optionJson string) error {
	history := &db_backend.History{
		MessageID:   messageId,
		UserID:      userId,
		CommandName: commandName,
		OptionJson:  optionJson,
		Created:     time.Now().Format("2006-01-02 15:04:05"),
	}
	err := ucs.Db.Db.Create(history).Error
	return err
}

// 写入图片信息
func (ucs *UserCenterService) WriteUserHistoryImages(messageId string, userId string, images string) error {
	err := ucs.Db.Db.Model(&db_backend.History{}).Where("message_id = ? AND user_id = ?", messageId, userId).Update("images", images).Error
	return err
}

// 获取历史记录
func (ucs *UserCenterService) GetUserHistoryOptWithMessageId(messageId string, commandName string) (string, error) {
	history := &db_backend.History{}
	err := ucs.Db.Db.Where("message_id = ? AND command_name = ?", messageId, commandName).First(history).Error
	if err != nil {
		return "", err
	}
	return history.OptionJson, nil
}

// 获取用户总数
func (ucs *UserCenterService) GetUserCount() (int64, error) {
	var count int64
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Count(&count).Error
	return count, err
}

// 判断当前用户是否为管理员
func (ucs *UserCenterService) IsAdmin(id string) (bool, error) {
	userInfo, err := ucs.GetUserInfo(id)
	if err != nil {
		return false, err
	}
	if userInfo == nil {
		return false, nil
	}
	roles := strings.Split(userInfo.Roles, ",")
	for _, role := range roles {
		if role == "admin" {
			return true, nil
		}
	}
	return false, nil
}

//

// 获取用户列表 判断当前用户是否为管理员，如果是管理员则返回所有用户，如果不是管理员则返回当前用户
func (ucs *UserCenterService) GetUserList(id string) ([]*UserInfo, error) {
	isAdmin, err := ucs.IsAdmin(id)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		dbUserInfos := []*db_backend.UserInfo{}
		err := ucs.Db.Db.Find(&dbUserInfos).Error
		if err != nil {
			return nil, err
		}
		return func() []*UserInfo {
			var userInfos []*UserInfo
			for _, v := range dbUserInfos {
				stableConfig := StableConfig{}
				if v.StableConfig != "{}" {
					err := json.NewDecoder(strings.NewReader(v.StableConfig)).Decode(&stableConfig)
					if err != nil {
						return nil
					}
				}
				userInfos = append(userInfos, &UserInfo{
					Enable:       v.Enable,
					Name:         v.Name,
					Id:           v.ID,
					Roles:        v.Roles,
					StableConfig: stableConfig,
				})

			}
			return userInfos
		}(), nil
	} else {
		userInfo, err := ucs.GetUserInfo(id)
		if err != nil {
			return nil, err
		}
		if userInfo == nil {
			return nil, nil
		}
		return []*UserInfo{userInfo}, nil
	}
}
