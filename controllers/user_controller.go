package controllers

import (
	"aiquiz/models/dto"
	"aiquiz/services"
	"aiquiz/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// ListUsers 条件分页查询所有用户
func (u *UserController) ListUsers(c *gin.Context) {
	// 验证是否为管理员
	role, exists := c.Get("role")
	if !exists || role != "admin" {
		utils.FailMsg(c, utils.ERROR_UNAUTHORIZED, "无权限访问")
		return
	}
	var req dto.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ParamError(c)
		return
	}
	// 默认分页大小
	req.Page = utils.NewPage(req.PageNum, req.PageSize)

	// 分页查询所有用户
	users, total, err := u.UserService.ListUsersByPage(c.Request.Context(), &req)
	if err != nil {
		utils.ServerErrorWithMsg(c, "获取用户列表失败"+err.Error())
		return
	}

	// 转为响应数据
	var usersResponse = make([]dto.UsersRes, 0, len(users))
	for _, user := range users {
		userRes := dto.UsersRes{
			ID:        user.ID,
			Username:  user.Username,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		usersResponse = append(usersResponse, userRes)
	}
	utils.SuccessMsg(c, utils.NewPageResult(usersResponse, total, req.PageNum, req.PageSize), "")
}

// UpdateUser 更新用户信息
func (u *UserController) UpdateUser(c *gin.Context) {
	// 获取路径中的userID并转换
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.ParamError(c)
		return
	}
	// 校验是否本人或管理员操作
	if userID != c.GetInt("user_id") && c.GetString("role") != "admin" {
		utils.FailMsg(c, utils.ERROR_UNAUTHORIZED, "无操作权限")
		return
	}
	var updateUserReq dto.UpdateUserReq
	if err := c.ShouldBindJSON(&updateUserReq); err != nil {
		utils.ParamError(c)
		return
	}
	// 更新用户信息
	err = u.UserService.UpdateUser(c.Request.Context(), userID, updateUserReq.Username, updateUserReq.Password)
	if err != nil {
		utils.ServerErrorWithMsg(c, "更新用户信息失败")
		return
	}
	utils.SuccessMsg(c, nil, "更新用户信息成功")
}

// DeleteUser 删除用户以及其级联数据
func (u *UserController) DeleteUser(c *gin.Context) {
	userID := c.GetInt("user_id")
	role, _ := c.Get("role")
	deleteUserID := c.Param("id")
	deleteUserIDInt, err := strconv.Atoi(deleteUserID)
	if err != nil {
		utils.ParamError(c)
	}
	// 校验是否本人或管理员操作
	if deleteUserIDInt != userID && role != "admin" {
		utils.FailMsg(c, utils.ERROR_NOT_PERMISSION, "无操作权限")
		return
	}
	// 删除数据
	err = u.UserService.DeleteUser(c.Request.Context(), deleteUserIDInt)
	if err != nil {
		utils.ServerErrorWithMsg(c, "删除用户失败")
		return
	}
	utils.Ok(c)
}
