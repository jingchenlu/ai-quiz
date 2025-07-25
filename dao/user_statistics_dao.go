package dao

import (
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"context"
	"gorm.io/gorm"
	"time"
)

type UserStatisticsDao struct {
	DB *gorm.DB
}

func NewUserStatisticsDao(db *gorm.DB) *UserStatisticsDao {
	return &UserStatisticsDao{DB: db}
}

// GetQuestionCount 统计用户出题次数
func (dao *UserStatisticsDao) GetQuestionCount(c context.Context, userID int) (int, error) {
	var count int64
	err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

// GetPaperCount 统计用户试卷数量
func (dao *UserStatisticsDao) GetPaperCount(c context.Context, userID int) (int, error) {
	var count int64
	err := dao.DB.WithContext(c).
		Model(&model.Paper{}).
		Where("creator_id = ?", userID).
		Count(&count).Error
	return int(count), err
}

// GetQuestionTypeDistribution 统计题目类型分布
func (dao *UserStatisticsDao) GetQuestionTypeDistribution(c context.Context, userID int) ([]dto.TypeDistribution, error) {
	var distribution []dto.TypeDistribution
	err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Where("user_id = ?", userID).
		Select("question_type as type, count(*) as count").
		Group("question_type").
		Scan(&distribution).Error
	return distribution, err
}

// GetLanguageTypeDistribution 统计语言类型分布
func (dao *UserStatisticsDao) GetLanguageTypeDistribution(c context.Context, userID int) ([]dto.LanguageDistribution, error) {
	var distribution []dto.LanguageDistribution
	err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Where("user_id = ?", userID).
		Select("language , count(*) as count").
		Group("language").
		Scan(&distribution).Error
	return distribution, err
}

// GetActiveTimeData 获取用户活跃时间原始数据
func (dao *UserStatisticsDao) GetActiveTimeData(c context.Context, userID int, startTime time.Time) ([]time.Time, error) {
	// 查询指定时间范围内的题目创建时间
	var questionTimes []time.Time
	if err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Select("created_at").
		Where("user_id = ? AND created_at >= ?", userID, startTime).
		Find(&questionTimes).Error; err != nil {
		return nil, err
	}

	// 查询指定时间范围内的试卷创建时间
	var paperTimes []time.Time
	if err := dao.DB.WithContext(c).
		Model(&model.Paper{}).
		Select("created_at").
		Where("creator_id = ? AND created_at >= ?", userID, startTime).
		Find(&paperTimes).Error; err != nil {
		return nil, err
	}

	// 合并时间数据
	return append(questionTimes, paperTimes...), nil
}
