package model

import (
	"gorm.io/gorm"
	"time"
)

// Paper 试卷模型
type Paper struct {
	ID          int            `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description" gorm:"type:text"`
	TotalScore  int            `json:"total_score" gorm:"default:100"`
	CreatorID   int            `json:"creator_id" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// 关联
	Creator *User `json:"creator" gorm:"foreignKey:CreatorID"`
	// PaperQuestion 表中的 PaperID 字段是关联到 Paper 表的外键。
	Questions []PaperQuestion `json:"questions" gorm:"foreignKey:PaperID"`
}

func (Paper) TableName() string {
	return "papers"
}
