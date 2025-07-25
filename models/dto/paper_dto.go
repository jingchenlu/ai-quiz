package dto

import "aiquiz/utils"

type GeneratePaperReq struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	TotalScore  int    `json:"total_score"`
}

type UpdatePaperReq struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TotalScore  int    `json:"total_score"`
}

// PaperListReq 条件分页查询试卷请求参数
type PaperListReq struct {
	utils.Page
	Title       string `form:"title"`
	Description string `form:"description"`
}

// AddPaperQuestionsReq 向试卷添加题目请求参数
type AddPaperQuestionsReq struct {
	QuestionID int `json:"question_id" validate:"required"`
	Score      int `json:"score"`
}

// PaperListRes 返回列表时的试卷结构
type PaperListRes struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	TotalScore  int    `json:"total_score"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	UserID      int    `json:"user_id"`
	UserName    string `json:"username"`
}

// PaperDetailRes 获取试卷详情结构体
type PaperDetailRes struct {
	PaperListRes
	Questions []QuestionRes `json:"questions"`
}

type QuestionOrderReq struct {
	QuestionID    int `json:"question_id" validate:"required"`
	QuestionOrder int `json:"question_order" validate:"required"`
	Score         int `json:"score"`
}
