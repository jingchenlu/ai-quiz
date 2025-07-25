package model

import (
	"gorm.io/gorm"
	"time"
)

// PaperQuestion 试卷题目关联模型
type PaperQuestion struct {
	ID            int            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	PaperID       int            `json:"paper_id" gorm:"not null"`
	QuestionID    int            `json:"question_id" gorm:"not null"`
	QuestionOrder int            `json:"question_order" gorm:"not null"` // 题目顺序
	Score         int            `json:"score" gorm:"default:5"`         // 该题分值
	CreatedAt     time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联
	Question *Question `json:"question" gorm:"foreignKey:QuestionID"`
}

func (PaperQuestion) TableName() string {
	return "paper_questions" // 或自定义的表名
}
