package models

import "gorm.io/gorm"

// golang的实体类 - 也就是模型 model
type User struct {
	// id等基础字段
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}
