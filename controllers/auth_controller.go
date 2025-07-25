package controllers

import (
	"aiquiz/config"
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"aiquiz/services"
	"aiquiz/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

type AuthController struct {
	UserService *services.UserService
}

func NewAuthController(service *services.UserService) *AuthController {
	return &AuthController{
		UserService: service,
	}
}

// Register 用户注册
func (a *AuthController) Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}
	// 创建用户
	err := a.UserService.Create(c.Request.Context(), &model.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	})
	if err != nil {
		// 校验是否唯一键冲突
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.FailMsg(c, utils.ERROR_USER_ALRAEDY_EXSIT, "用户名已存在")
			return
		}
		utils.ServerErrorWithMsg(c, "用户创建失败"+err.Error())
		return
	}
	utils.Ok(c)
}

// Login 用户登录
func (a *AuthController) Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}
	// 验证登录
	user, err := a.UserService.Login(c.Request.Context(), &model.User{
		Username:     req.Username,
		PasswordHash: req.Password,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.FailMsg(c, utils.ERROR_USER_OR_PASSWORD_NOT_EXIST, "用户名不存在")
			return
		}
		utils.FailMsg(c, utils.ERROR_USER_OR_PASSWORD_NOT_EXIST, "登录失败"+err.Error())
		return
	}
	// 生成token
	token, err := utils.GenerateToken(user)
	if err != nil {
		utils.ServerErrorWithMsg(c, "token生成失败"+err.Error())
		return
	}
	// 设置cookie以及到期时间
	jwtConfig := config.GetJWTConfig()
	expiresIn := int(jwtConfig.TokenExpiry.Seconds())
	c.SetCookie("token", token, expiresIn, "/", "", false, true)
	// 返回响应
	utils.SuccessMsg(c, nil, "登录成功")
}
