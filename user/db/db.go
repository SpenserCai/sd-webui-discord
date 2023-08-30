/*
 * @Author: SpenserCai
 * @Date: 2023-08-30 21:21:40
 * @version:
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-08-30 21:30:44
 * @Description: file content
 */
package db

import "gorm.io/gorm"

type DbotDb struct {
	Db *gorm.DB
}

// TODO: init db with config
