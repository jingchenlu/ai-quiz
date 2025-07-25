package dao

import (
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"context"
	"gorm.io/gorm"
	"time"
)

type SystemStatisticsDao struct {
	DB *gorm.DB
}

func NewSystemStatisticsDao(db *gorm.DB) *SystemStatisticsDao {
	return &SystemStatisticsDao{DB: db}
}

// GetTotalUserCount 获取总用户数
func (dao *SystemStatisticsDao) GetTotalUserCount(c context.Context) (int, error) {
	var count int64
	err := dao.DB.WithContext(c).
		Model(&model.User{}).
		Count(&count).Error
	return int(count), err
}

// GetTotalQuestionCount 获取总题目数
func (dao *SystemStatisticsDao) GetTotalQuestionCount(c context.Context) (int, error) {
	var count int64
	err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Count(&count).Error
	return int(count), err
}

// GetTotalPaperCount 获取总试卷数
func (dao *SystemStatisticsDao) GetTotalPaperCount(c context.Context) (int, error) {
	var count int64
	err := dao.DB.WithContext(c).
		Model(&model.Paper{}).
		Count(&count).Error
	return int(count), err
}

// GetLanguageDistribution 获取编程语言分布
func (dao *SystemStatisticsDao) GetLanguageDistribution(c context.Context) ([]dto.LanguageDistribution, error) {
	var distribution []dto.LanguageDistribution
	err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Select("language, count(*) as count").
		Group("language").
		Scan(&distribution).Error
	return distribution, err
}

// GetAIModelUsage 获取AI模型使用情况
func (dao *SystemStatisticsDao) GetAIModelUsage(c context.Context) ([]dto.AIModelDistribution, error) {
	var usage []dto.AIModelDistribution
	err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Select("ai_model as model_name, count(*) as count").
		Where("ai_model != ''").
		Group("ai_model").
		Scan(&usage).Error
	return usage, err
}

// GetPaperQuestionDistribution 获取试卷题目数量分布
func (dao *SystemStatisticsDao) GetPaperQuestionDistribution(c context.Context) ([]dto.PaperQuestionDistribution, error) {
	// 先查询每个试卷的题目数量
	var paperQuestionCounts []struct {
		PaperID int
		Count   int
	}
	err := dao.DB.WithContext(c).
		Model(&model.PaperQuestion{}).
		Select("paper_id, count(question_id) as count").
		Group("paper_id").
		Scan(&paperQuestionCounts).Error
	if err != nil {
		return nil, err
	}

	// 按范围统计
	distributionMap := map[string]int{
		"1-10":  0,
		"11-20": 0,
		"21-50": 0,
		"51+":   0,
	}

	for _, item := range paperQuestionCounts {
		switch {
		case item.Count <= 10:
			distributionMap["1-10"]++
		case item.Count <= 20:
			distributionMap["11-20"]++
		case item.Count <= 50:
			distributionMap["21-50"]++
		default:
			distributionMap["51+"]++
		}
	}

	// 转换为DTO
	result := make([]dto.PaperQuestionDistribution, 0, len(distributionMap))
	for r, c := range distributionMap {
		result = append(result, dto.PaperQuestionDistribution{
			Range: r,
			Count: c,
		})
	}
	return result, nil
}

// GetSystemActivityData 获取系统活跃度原始数据
func (dao *SystemStatisticsDao) GetSystemActivityData(c context.Context, startTime time.Time) ([]time.Time, error) {
	// 用户注册时间
	var userTimes []time.Time
	if err := dao.DB.WithContext(c).
		Model(&model.User{}).
		Select("created_at").
		Where("created_at >= ?", startTime).
		Find(&userTimes).Error; err != nil {
		return nil, err
	}

	// 题目创建时间
	var questionTimes []time.Time
	if err := dao.DB.WithContext(c).
		Model(&model.Question{}).
		Select("created_at").
		Where("created_at >= ?", startTime).
		Find(&questionTimes).Error; err != nil {
		return nil, err
	}

	// 试卷创建时间
	var paperTimes []time.Time
	if err := dao.DB.WithContext(c).
		Model(&model.Paper{}).
		Select("created_at").
		Where("created_at >= ?", startTime).
		Find(&paperTimes).Error; err != nil {
		return nil, err
	}

	// 合并所有时间
	return append(append(userTimes, questionTimes...), paperTimes...), nil
}
