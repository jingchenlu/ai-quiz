package middlewares

import (
	"aiquiz/dao/model"
	"aiquiz/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

func PaperAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取路径参数
		paperIDStr := c.Param("paper_id")
		paperID, err := strconv.Atoi(paperIDStr)
		if err != nil || paperID <= 0 {
			utils.FailMsg(c, utils.BAD_REQUEST, "无效的 paper_id")
			c.Abort()
			return
		}
		userID := c.GetInt("user_id")
		role := c.GetString("role")
		if role != "admin" {
			// 检查该用户的试卷权限
			var paper model.Paper
			err = db.WithContext(c.Request.Context()).Model(&model.Paper{}).
				Where("id = ? AND creator_id = ?", paperID, userID).
				Take(&paper).Error
			if err != nil || paper.ID == 0 {
				utils.NotPermission(c)
				c.Abort()
				return
			}
		}
		c.Set("paper_id", paperID)
		c.Next()
	}
}
