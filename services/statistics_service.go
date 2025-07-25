package services

import (
	"aiquiz/dao"
	"aiquiz/models/dto"
	"context"
	"fmt"
	"time"
)

type StatisticsService struct {
	userDao             *dao.UserDao
	userStatisticsDao   *dao.UserStatisticsDao
	systemStatisticsDao *dao.SystemStatisticsDao
}

func NewStatisticService(
	userDao *dao.UserDao,
	userStatisticsDao *dao.UserStatisticsDao,
	systemStatisticsDao *dao.SystemStatisticsDao,
) *StatisticsService {
	return &StatisticsService{
		userDao:             userDao,
		userStatisticsDao:   userStatisticsDao,
		systemStatisticsDao: systemStatisticsDao,
	}
}

func (s *StatisticsService) GetUserStatistics(c context.Context, userID int) (*dto.UserStatisticsRes, error) {
	// 获取用户基本信息
	user, err := s.userDao.GetUserByID(c, userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}
	// 统计出题次数
	questionCount, err := s.userStatisticsDao.GetQuestionCount(c, userID)
	if err != nil {
		return nil, fmt.Errorf("统计出题次数失败: %v", err)
	}

	// 统计试卷数量
	paperCount, err := s.userStatisticsDao.GetPaperCount(c, userID)
	if err != nil {
		return nil, fmt.Errorf("统计试卷数量失败: %v", err)
	}

	// 统计题目类型分布
	typeDistribution, err := s.userStatisticsDao.GetQuestionTypeDistribution(c, userID)
	if err != nil {
		return nil, fmt.Errorf("统计题目类型分布失败: %v", err)
	}

	// 统计语言分布
	languageDistribution, err := s.userStatisticsDao.GetLanguageTypeDistribution(c, userID)
	if err != nil {
		return nil, fmt.Errorf("统计语言分布失败: %v", err)
	}

	// 分析活跃时间 (获取最近一年的数据)
	oneYearAgo := time.Now().AddDate(-1, 0, 0)
	timeData, err := s.userStatisticsDao.GetActiveTimeData(c, userID, oneYearAgo)
	if err != nil {
		return nil, fmt.Errorf("获取活跃时间数据失败: %v", err)
	}

	activeTimeAnalysis := s.analyzeActiveTime(timeData)

	// 组装结果
	return &dto.UserStatisticsRes{
		UserID:                   userID,
		Username:                 user.Username,
		QuestionCount:            questionCount,
		PaperCount:               paperCount,
		QuestionTypeDistribution: typeDistribution,
		LanguageDistribution:     languageDistribution,
		ActiveTimeAnalysis:       activeTimeAnalysis,
	}, nil
}

// analyzeActiveTime 分析活跃时间数据
func (s *StatisticsService) analyzeActiveTime(times []time.Time) dto.ActiveTimeAnalysis {
	analysis := dto.ActiveTimeAnalysis{
		Daily:   make(map[string]int),
		Monthly: make(map[string]int),
	}

	// 最近30天的日期范围
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	for _, t := range times {
		// 按小时统计
		hour := t.Hour()
		analysis.Hourly[hour]++

		// 按天统计(只统计最近30天)
		if t.After(thirtyDaysAgo) {
			dayKey := t.Format("2006-01-02")
			analysis.Daily[dayKey]++
		}

		// 按月统计
		monthKey := t.Format("2006-01")
		analysis.Monthly[monthKey]++
	}

	return analysis
}

func (s *StatisticsService) GetSystemStatistics(c context.Context) (*dto.SystemStatisticsRes, error) {
	// 总用户数
	totalUser, err := s.systemStatisticsDao.GetTotalUserCount(c)
	if err != nil {
		return nil, fmt.Errorf("获取总用户数失败: %v", err)
	}

	// 总题目数
	totalQuestion, err := s.systemStatisticsDao.GetTotalQuestionCount(c)
	if err != nil {
		return nil, fmt.Errorf("获取总题目数失败: %v", err)
	}

	// 总试卷数
	totalPaper, err := s.systemStatisticsDao.GetTotalPaperCount(c)
	if err != nil {
		return nil, fmt.Errorf("获取总试卷数失败: %v", err)
	}

	// 编程语言分布
	languageDist, err := s.systemStatisticsDao.GetLanguageDistribution(c)
	if err != nil {
		return nil, fmt.Errorf("获取语言分布失败: %v", err)
	}

	// AI模型使用情况
	aiUsage, err := s.systemStatisticsDao.GetAIModelUsage(c)
	if err != nil {
		return nil, fmt.Errorf("获取AI模型使用情况失败: %v", err)
	}

	// 试卷题目数量分布
	paperQuestionDist, err := s.systemStatisticsDao.GetPaperQuestionDistribution(c)
	if err != nil {
		return nil, fmt.Errorf("获取试卷题目分布失败: %v", err)
	}

	// 活跃度分析,通过user,paper,questions的创建时间来分析（最近90天）
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	activityData, err := s.systemStatisticsDao.GetSystemActivityData(c, ninetyDaysAgo)
	if err != nil {
		return nil, fmt.Errorf("获取活跃度数据失败: %v", err)
	}

	activityAnalysis := s.analyzeSystemActivity(activityData)

	return &dto.SystemStatisticsRes{
		TotalUserCount:            totalUser,
		TotalQuestionCount:        totalQuestion,
		TotalPaperCount:           totalPaper,
		LanguageDistribution:      languageDist,
		AIModelUsage:              aiUsage,
		PaperQuestionDistribution: paperQuestionDist,
		ActivityAnalysis:          activityAnalysis,
	}, nil
}

// analyzeSystemActivity 分析系统活跃度
func (s *StatisticsService) analyzeSystemActivity(times []time.Time) dto.SystemActivityAnalysis {
	analysis := dto.SystemActivityAnalysis{
		Daily:  make(map[string]int),
		Weekly: make(map[string]int),
	}

	for _, t := range times {
		// 按天统计
		dayKey := t.Format("2006-01-02")
		analysis.Daily[dayKey]++

		// 按周统计（格式：2006-W01）
		year, week := t.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)
		analysis.Weekly[weekKey]++
	}

	return analysis
}
