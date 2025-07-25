package controllers

import (
	"aiquiz/services"
	"aiquiz/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type StatisticController struct {
	StatisticsService *services.StatisticsService
}

func NewStatisticController(statisticsService *services.StatisticsService) *StatisticController {
	return &StatisticController{
		StatisticsService: statisticsService,
	}
}

// GetUserStatistics 获取用户统计信息
func (u *StatisticController) GetUserStatistics(c *gin.Context) {
	UserIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(UserIDStr)
	if err != nil {
		utils.BadRequestWithMsg(c, "用户ID必须为整型")
	}
	statistics, err := u.StatisticsService.GetUserStatistics(c.Request.Context(), userID)
	if err != nil {
		utils.ServerErrorWithMsg(c, "统计失败"+err.Error())
		return
	}
	utils.SuccessMsg(c, statistics, "")
}

// GetSystemStatistics 获取系统整体统计
func (u *StatisticController) GetSystemStatistics(c *gin.Context) {
	statistics, err := u.StatisticsService.GetSystemStatistics(c.Request.Context())
	if err != nil {
		utils.ServerErrorWithMsg(c, "系统统计失败: "+err.Error())
		return
	}
	utils.SuccessMsg(c, statistics, "系统统计获取成功")
}
