package dto

import "aiquiz/utils"

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginReq 登录请求
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserListReq struct {
	Username string `form:"username"`
	utils.Page
}

type UpdateUserReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UsersRes 用户信息响应
type UsersRes struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
