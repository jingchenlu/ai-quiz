package middlewares

import (
	"aiquiz/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT登录认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从cookie中取出token
		token, err := c.Cookie("token")
		if err != nil {
			// 再尝试从header中取出token
			token = c.GetHeader("Authorization")
			if token == "" {
				utils.FailMsg(c, utils.ERROR_UNAUTHORIZED, "未提供授权令牌")
				c.Abort()
				return
			}
		}
		claims, err := utils.ValidateToken(token)
		if err != nil {
			utils.FailMsg(c, utils.ERROR_UNAUTHORIZED, "授权令牌无效")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")
		if !ok {
			utils.FailMsg(c, utils.ERROR_UNAUTHORIZED, "未提供授权令牌")
			c.Abort()
			return
		}
		if role != "admin" {
			utils.FailMsg(c, utils.ERROR_NOT_PERMISSION, "无管理员权限")
		}
	}
}
