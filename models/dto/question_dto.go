package dto

import (
	"aiquiz/utils"
	"aiquiz/utils/enums"
)

// Option 题目选项结构体,value依次为2的次幂，便于用移位&进行少选错选的判断
type Option struct {
	Content string `json:"content"`
	Value   int    `json:"value"`
}

// Question 单道题目结构体(用于解析ai模型生成的题目)
type Question struct {
	Title       string   `json:"title"`
	Options     []Option `json:"options"`
	Answer      int      `json:"answer"`
	Explanation string   `json:"explanation"`
}

// GenerateQuestionReq 生成题目请求结构体
type GenerateQuestionReq struct {
	Language     string             `json:"language" validate:"required"`
	QuestionType enums.QuestionType `json:"question_type" validate:"required"`
	Keywords     string             `json:"keywords" validate:"required"`
	Count        int                `json:"count" validate:"required"`
	AiModel      enums.AiModel      `json:"ai_model" validate:"required"`
}

// ConfirmQuestionReq 确认题目请求结构体
type ConfirmQuestionReq struct {
	Title        string             `json:"title" validate:"required"`
	Options      []Option           `json:"options" validate:"required"`
	Answer       int                `json:"answer" validate:"required"`
	Explanation  string             `json:"explanation" validate:"required"`
	QuestionType enums.QuestionType `json:"question_type" validate:"required"`
	Language     string             `json:"language" validate:"required"`
	AiModel      enums.AiModel      `json:"ai_model" validate:"required"`
	Keywords     string             `json:"keywords" validate:"required"`
}

// ListQuestionsReq 分页获取题目列表（根据条件选择）
type ListQuestionsReq struct {
	utils.Page
	Title        string             `form:"title"`
	QuestionType enums.QuestionType `form:"question_type"`
	Language     string             `form:"language"`
	AiModel      enums.AiModel      `form:"ai_model"`
	Keywords     string             `form:"keywords"`
}

type UpdateQuestionReq struct {
	Question
	QuestionType enums.QuestionType `json:"question_type"`
	Language     string             `json:"language"`
	Keywords     string             `json:"keywords"`
}

// QuestionRes 查询题目列表返回结构体
type QuestionRes struct {
	ID int `json:"id"`
	Question
	QuestionType string `json:"question_type"`
	Language     string `json:"language"`
	Keywords     string `json:"keywords"`
	AiModel      string `json:"ai_model"`
	CreateAt     string `json:"created_at"`
	UserName     string `json:"username"`
	UserID       int    `json:"user_id"`
}

// GenerateQuestionRes 生成题目返回结构体
type GenerateQuestionRes struct {
	Question
	QuestionType string `json:"question_type"`
	Language     string `json:"language"`
	Keywords     string `json:"keywords"`
	AiModel      string `json:"ai_model"`
}
