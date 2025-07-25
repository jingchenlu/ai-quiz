package model

import (
	"gorm.io/gorm"
	"time"
)

// Question 题目模型
type Question struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Title        string         `json:"title" gorm:"type:text;not null"`
	QuestionType string         `json:"question_type" gorm:"size:20;not null"` // 'single' 或 'multiple'
	Options      string         `json:"options" gorm:"type:text;not null"`     // JSON格式存储选项
	Answer       string         `json:"answer" gorm:"type:text;not null"`
	Explanation  string         `json:"explanation" gorm:"type:text"`
	Keywords     string         `json:"keywords" gorm:"size:255"`
	Language     string         `json:"language" gorm:"size:50;not null"` // 编程语言
	AiModel      string         `json:"ai_model" gorm:"size:50;not null"` // 使用的AI模型
	UserID       int            `json:"user_id" gorm:"not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"` // 使用指针表示可为空

	// 关联
	User *User `json:"user" gorm:"foreignKey:UserID"`
}

func (Question) TableName() string {
	return "questions"
}
