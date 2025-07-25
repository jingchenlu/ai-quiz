package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`    // 状态码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 业务数据(可选)
}

// 状态码常量
const (
	SUCCESS                          = 200  // 成功
	BAD_REQUEST                      = 400  // 客户端错误
	ERROR_SERVER                     = 500  // 服务器错误
	ERROR_UNAUTHORIZED               = 401  // 未授权
	ERROR_NOTFOUND                   = 404  // 资源不存在
	ERROR_PARAM                      = 4001 // 请求参数有误
	ERROR_USER_ALRAEDY_EXSIT         = 4002 // 用户名已存在
	ERROR_USER_OR_PASSWORD_NOT_EXIST = 4003 // 用户名或密码错误
	ERROR_AI_GENERATE                = 4004 // AI生成错误
	ERROR_NOT_PERMISSION             = 4005 // 权限不足
	ERROR_RECORD_NOT_EXIST           = 4006
)

// SuccessMsg Success 成功响应
func SuccessMsg(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "操作成功"
	}
	c.JSON(http.StatusOK, Response{
		Code:    SUCCESS,
		Message: message,
		Data:    data,
	})
}

// FailMsg Fail 错误响应
func FailMsg(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// Ok 快捷方法：常见场景
func Ok(c *gin.Context) {
	SuccessMsg(c, nil, "操作成功")
}

func BadRequestWithMsg(c *gin.Context, message string) {
	FailMsg(c, BAD_REQUEST, message)
}

func ServerErrorWithMsg(c *gin.Context, message string) {
	FailMsg(c, ERROR_SERVER, message)
}
func ParamError(c *gin.Context) {
	FailMsg(c, ERROR_PARAM, "请求参数有误")
}
func NotPermission(c *gin.Context) {
	FailMsg(c, ERROR_NOT_PERMISSION, "无权限操作不属于您的资源")
}
