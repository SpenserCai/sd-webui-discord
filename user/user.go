/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 20:38:24
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-19 14:20:09
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

type UserInfo struct {
	Enable             bool                `json:"enable"`
	IsPrivate          bool                `json:"is_private"`
	Avatar             string              `json:"avatar"`
	Name               string              `json:"name"`
	Id                 string              `json:"id"`
	Roles              string              `json:"roles"`
	Created            string              `json:"created"`
	StableConfig       config.StableConfig `json:"stable_config"`
	CycleCredit        int                 `json:"cycle_credit"`         // 周期credit
	CreditUpdateCycle  string              `json:"credit_update_cycle"`  // credit更新周期 格式 1|H 1|D 1|W 1|M
	CycleCreditUpdated string              `json:"cycle_credit_updated"` // 周期credit更新时间，用于判断是否需要更新credit
	PlusCredit         int                 `json:"plus_credit"`          // 充值的credit
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

	stableConfig := config.StableConfig{}
	if userInfo.StableConfig != "" {
		err := json.NewDecoder(strings.NewReader(userInfo.StableConfig)).Decode(&stableConfig)
		if err != nil {
			return nil, err
		}
	}

	return &UserInfo{
		Enable:       userInfo.Enable,
		IsPrivate:    userInfo.IsPrivate,
		Avatar:       userInfo.Avatar,
		Name:         userInfo.Name,
		Id:           userInfo.ID,
		Roles:        userInfo.Roles,
		StableConfig: stableConfig,
		Created:      userInfo.Created,
	}, nil
}

func (ucs *UserCenterService) GetUserImageTotal(id string) (int64, error) {
	var count int64
	err := ucs.Db.Db.Model(&db_backend.History{}).Where("user_id = ? AND deleted = ?", id, false).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (ucs *UserCenterService) GetUsersImageTotal(ids []string) (map[string]int, error) {
	// 根据ids获取用户生成的图片数量，用一条sql语句获取
	infoMap := make(map[string]int)
	type UsersImageTotal struct {
		UserID string `gorm:"column:user_id"`
		Count  int    `gorm:"column:count"`
	}
	var usersImageTotal []UsersImageTotal
	err := ucs.Db.Db.Model(&db_backend.History{}).Select("user_id,count(*) as count").Where("user_id in (?) AND deleted = ?", ids, false).Group("user_id").Scan(&usersImageTotal).Error
	if err != nil {
		return nil, err
	}
	for _, v := range usersImageTotal {
		infoMap[v.UserID] = v.Count
	}
	return infoMap, nil
}

func (ucs *UserCenterService) GetUserInfoList(ids []string) ([]*UserInfo, error) {
	userInfos := []*UserInfo{}
	dbUserInfos := []*db_backend.UserInfo{}
	err := ucs.Db.Db.Where("id in (?)", ids).Find(&dbUserInfos).Error
	if err != nil {
		return nil, err
	}
	for _, v := range dbUserInfos {
		stableConfig := config.StableConfig{}
		if v.StableConfig != "{}" {
			err := json.NewDecoder(strings.NewReader(v.StableConfig)).Decode(&stableConfig)
			if err != nil {
				return nil, err
			}
		}
		userInfos = append(userInfos, &UserInfo{
			Enable:       v.Enable,
			IsPrivate:    v.IsPrivate,
			Avatar:       v.Avatar,
			Name:         v.Name,
			Id:           v.ID,
			Roles:        v.Roles,
			StableConfig: stableConfig,
			Created:      v.Created,
		})

	}
	return userInfos, nil
}

func (ucs *UserCenterService) SearchUserInfoList(page int, pageSize int, query map[string]interface{}) ([]*UserInfo, int, error) {
	// 判断是否有username字段，如果有则按照username进行模糊查询
	queryString := ""
	queryValues := []interface{}{}
	if query["username"] != nil {
		queryString += "name LIKE ? AND "
		queryValues = append(queryValues, "%"+query["username"].(string)+"%")
	}
	if query["id"] != nil {
		queryString += "id = ? AND "
		queryValues = append(queryValues, query["id"].(string))
	}
	if query["enable"] != nil {
		queryString += "enable = ? AND "
		queryValues = append(queryValues, query["enable"].(bool))
	}
	// 去掉最后一个AND
	queryString = func() string {
		if queryString == "" {
			return ""
		} else {
			return queryString[:len(queryString)-4]
		}
	}()
	dbUserInfos := []*db_backend.UserInfo{}
	err := ucs.Db.Db.Where(queryString, queryValues...).Offset((page - 1) * pageSize).Limit(pageSize).Find(&dbUserInfos).Error
	if err != nil {
		return nil, 0, err
	}
	// 获取总数
	var count int64
	err = ucs.Db.Db.Model(&db_backend.UserInfo{}).Where(queryString, queryValues...).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	userInfos := []*UserInfo{}
	for _, v := range dbUserInfos {
		stableConfig := config.StableConfig{}
		if v.StableConfig != "{}" {
			err := json.NewDecoder(strings.NewReader(v.StableConfig)).Decode(&stableConfig)
			if err != nil {
				return nil, 0, err
			}
		}
		userInfos = append(userInfos, &UserInfo{
			Enable:       v.Enable,
			IsPrivate:    v.IsPrivate,
			Avatar:       v.Avatar,
			Name:         v.Name,
			Id:           v.ID,
			Roles:        v.Roles,
			StableConfig: stableConfig,
			Created:      v.Created,
		})

	}
	return userInfos, int(count), nil

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
		ID:        user.Id,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Created:   time.Now().Format("2006-01-02 15:04:05"),
		Enable:    true,
		IsPrivate: false,
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
		ID:        user.Id,
		Name:      user.Name,
		Avatar:    user.Avatar,
		Enable:    user.Enable,
		Roles:     user.Roles,
		IsPrivate: user.IsPrivate,
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

func (ucs *UserCenterService) SetUserEnable(id string, enable bool) error {
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", id).Update("enable", enable).Error
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
	isPrivate := func() bool {
		isPriva, err := ucs.IsPrivate(userId)
		if err != nil {
			return false
		}
		return isPriva
	}()
	history := &db_backend.History{
		MessageID:   messageId,
		UserID:      userId,
		CommandName: commandName,
		OptionJson:  optionJson,
		Created:     time.Now().Format("2006-01-02 15:04:05"),
		Deleted:     false,
		IsPrivate:   isPrivate,
	}
	err := ucs.Db.Db.Create(history).Error
	return err
}

// 写入图片信息
func (ucs *UserCenterService) WriteUserHistoryImages(messageId string, userId string, images string, imagesBlurHash string) error {
	err := ucs.Db.Db.Model(&db_backend.History{}).Where("message_id = ? AND user_id = ?", messageId, userId).Update("images", images).Update("image_blur_hashs", imagesBlurHash).Error
	return err
}

// 删除历史记录（软删除）
func (ucs *UserCenterService) DeleteUserHistory(messageId string, userId string) error {
	err := ucs.Db.Db.Model(&db_backend.History{}).Where("message_id = ? AND user_id = ?", messageId, userId).Update("deleted", true).Error
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

// 获取用户历史记录列表分页获取 TODO:加上参数：excludePrivate
func (ucs *UserCenterService) GetUserHistoryList(userId string, cmd string, page int, pageSize int, excludePrivate bool) ([]*db_backend.History, int, error) {
	historyList := []*db_backend.History{}
	var err error
	// 之获取没有软删除的历史记录，同时返回总数,按照created的倒序排列
	if userId == "" {
		args := []interface{}{cmd, false}
		whereString := "command_name = ? AND deleted = ?"
		if excludePrivate {
			whereString += " AND is_private = ?"
			args = append(args, false)
		}
		err = ucs.Db.Db.Where(whereString, args...).Order("created desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&historyList).Error
	} else {
		err = ucs.Db.Db.Where("user_id = ? AND command_name = ? AND deleted = ?", userId, cmd, false).Order("created desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&historyList).Error
	}
	if err != nil {
		return nil, 0, err
	}
	var count int64
	if userId == "" {
		err = ucs.Db.Db.Model(&db_backend.History{}).Where("command_name = ? AND deleted = ?", cmd, false).Count(&count).Error
	} else {
		err = ucs.Db.Db.Model(&db_backend.History{}).Where("user_id = ? AND command_name = ? AND deleted = ?", userId, cmd, false).Count(&count).Error
	}
	if err != nil {
		return nil, 0, err
	}
	return historyList, int(count), nil
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

// 获取用户是否是Private
func (ucs *UserCenterService) IsPrivate(id string) (bool, error) {
	// 根据id判断用户是否是private
	isPrivate := false
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", id).Select("is_private").Scan(&isPrivate).Error
	if err != nil {
		return false, err
	}
	return isPrivate, nil
}

// 设置用户是否是Private
func (ucs *UserCenterService) SetUserPrivate(id string, isPrivate bool) error {
	err := ucs.Db.Db.Model(&db_backend.UserInfo{}).Where("id = ?", id).Update("is_private", isPrivate).Error
	return err
}

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
				stableConfig := config.StableConfig{}
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
