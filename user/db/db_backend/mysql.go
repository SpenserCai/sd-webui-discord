/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 00:33:12
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 00:44:44
 * @Description: file content
 */
package db_backend

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func (backend DbBackend) CreateDbMysqlConnect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
