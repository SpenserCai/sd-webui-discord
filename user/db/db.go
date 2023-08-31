/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 21:21:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 00:58:05
 * @Description: file content
 */
package db

import (
	"errors"
	"reflect"
	"strings"

	"github.com/SpenserCai/sd-webui-discord/config"
	"github.com/SpenserCai/sd-webui-discord/user/db/db_backend"
	"gorm.io/gorm"
)

type BotDb struct {
	Db *gorm.DB
}

func NewBotDb(dbCfg *config.DbConfig) (*BotDb, error) {
	dbType := dbCfg.Type

	dbCreateName := "CreateDb" + strings.ToUpper(dbType[:1]) + dbType[1:] + "Connect"
	pkgValue := reflect.ValueOf(db_backend.DbBackend{})
	methodValue := pkgValue.MethodByName(dbCreateName)
	if !methodValue.IsValid() {
		return nil, errors.New("db type not support")
	}

	createDbFunc := methodValue.Interface().(func(string) (*gorm.DB, error))
	db, err := createDbFunc(dbCfg.DSN)
	if err != nil {
		return nil, err
	}

	botDb := &BotDb{
		Db: db,
	}
	return botDb, nil
}
