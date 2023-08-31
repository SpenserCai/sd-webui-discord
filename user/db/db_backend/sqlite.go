/*
 * @Author: SpenserCai
 * @Date: 2023-08-31 00:08:39
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-31 00:44:55
 * @Description: file content
 */
package db_backend

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (backend DbBackend) CreateDbSqliteConnect(dsn string) (*gorm.DB, error) {
	_, err := os.Stat(dsn)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return db, err
}
