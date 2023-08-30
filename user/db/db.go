/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 21:21:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-30 21:50:12
 * @Description: file content
 */
package db

import "gorm.io/gorm"

type BotDb struct {
	Db *gorm.DB
}

// TODO: init db with config
