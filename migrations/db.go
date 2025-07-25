package migrations

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

const (
	initSQLPath = "migrations/init.sql"
)

func InitDB(dbPath string) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Millisecond * 200,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	// 打开数据库连接, 不存在会自动创建db
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(fmt.Errorf("打开数据库失败: %v", err))
	}

	// 执行 init.sql
	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("获取底层数据库失败: %v", err))
	}

	sqlContent, err := os.ReadFile(initSQLPath)
	if err != nil {
		panic(fmt.Errorf("读取 init.sql 失败: %v", err))
	}

	if _, err := sqlDB.Exec(string(sqlContent)); err != nil {
		panic(fmt.Errorf("执行 init.sql 失败: %v", err))
	}

	return db
}
