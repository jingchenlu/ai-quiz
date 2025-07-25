package main

import (
	"aiquiz/config"
	"aiquiz/dao/model"
	"fmt"
	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func main() {
	appConfig := config.GetConfig(false)

	// 打开数据库连接, 不存在会自动创建db
	db, err := gorm.Open(sqlite.Open(appConfig.DBPath), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("打开数据库失败: %v", err))
	}
	// 根据定义的model建表
	err = db.AutoMigrate(
		&model.User{},
		&model.Question{},
		&model.Paper{},
		&model.PaperQuestion{},
	)
	if err != nil {
		panic(fmt.Errorf("建表失败: %v", err))
	}
	// 默认带一个admin用户，默认密码是123456
	adminUser := &model.User{Username: "admin", PasswordHash: "123456"}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminUser.PasswordHash), bcrypt.DefaultCost)
	adminUser.PasswordHash = string(hashedPassword)
	if err := db.Create(adminUser).Error; err != nil {
		log.Printf("创建默认用户失败: %v", err)
	}

	log.Println("数据库初始化完成")
}
