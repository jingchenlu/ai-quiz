package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Username     string         `json:"username" gorm:"size:50;unique;not null"`
	PasswordHash string         `json:"password_hash" gorm:"size:255;not null"` // 不序列化到JSON
	Role         string         `json:"role" gorm:"size:20;default:user"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"` // 使用指针表示可为空
}

// TableName 指定User结构体对应的表名
func (User) TableName() string {
	return "users"
}
