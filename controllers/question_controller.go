package controllers

import (
	"aiquiz/ai"
	"aiquiz/config"
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"aiquiz/services"
	"aiquiz/utils"
	"aiquiz/utils/enums"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

type QuestionController struct {
	QuestionService *services.QuestionService
}

func NewQuestionController(questionService *services.QuestionService) *QuestionController {
	return &QuestionController{QuestionService: questionService}
}

// GenerateQuestion 调用ai模型生成题目并验证
func (q *QuestionController) GenerateQuestion(c *gin.Context) {
	var req dto.GenerateQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	// 验证参数
	appConfig := config.GetConfig(true)
	if !enums.IsSupportedQuestionType(req.QuestionType) {
		utils.BadRequestWithMsg(c, "无效的题目类型，必须是 'single' 或 'multiple'")
		return
	}
	if req.Count < 1 || req.Count > 10 {
		utils.BadRequestWithMsg(c, "题目数量必须在1到10之间")
		return
	}
	if _, ok := appConfig.SupportedLanguages[req.Language]; !ok {
		utils.BadRequestWithMsg(c, "无效的语言")
		return
	}
	if !enums.IsSupportedAiModel(req.AiModel) {
		utils.BadRequestWithMsg(c, "无效的AI模型")
		return
	}

	// 调用ai模型
	generatedQuestions, err := ai.GenerateQuestions(string(req.AiModel), req.Language, string(req.QuestionType), req.Keywords, req.Count)
	// 应该重试
	if err != nil {
		utils.FailMsg(c, utils.ERROR_AI_GENERATE, "生成题目失败"+err.Error())
		return
	}
	if generatedQuestions == nil || len(generatedQuestions.Questions) == 0 {
		utils.FailMsg(c, utils.ERROR_AI_GENERATE, "生成题目失败,请重试")
		return
	}
	var questionResponseList []dto.GenerateQuestionRes
	for _, question := range generatedQuestions.Questions {
		questionRes := dto.GenerateQuestionRes{
			Question:     question,
			QuestionType: string(req.QuestionType),
			Language:     req.Language,
			AiModel:      string(req.AiModel),
			Keywords:     req.Keywords,
		}
		questionResponseList = append(questionResponseList, questionRes)
	}
	utils.SuccessMsg(c, questionResponseList, "生成题目成功")
}

// ConfirmQuestions 确认题目（入库）
func (q *QuestionController) ConfirmQuestions(c *gin.Context) {
	var reqs []dto.ConfirmQuestionReq

	if err := c.ShouldBindJSON(&reqs); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	// 验证参数
	appConfig := config.GetConfig(true)
	if reqs == nil || len(reqs) == 0 {
		utils.BadRequestWithMsg(c, "请提供题目")
		return
	}
	questions := make([]model.Question, 0, len(reqs))
	// 转换为模型
	for _, req := range reqs {
		if !enums.IsSupportedQuestionType(req.QuestionType) {
			utils.BadRequestWithMsg(c, "无效的题目类型，必须是 'single' 或 'multiple'")
			return
		}
		if _, ok := appConfig.SupportedLanguages[req.Language]; !ok {
			utils.BadRequestWithMsg(c, "无效的语言")
			return
		}
		if !enums.IsSupportedAiModel(req.AiModel) {
			utils.BadRequestWithMsg(c, "无效的AI模型")
			return
		}
		// 序列化 Options 为 JSON 字符串
		optionBytes, err := json.Marshal(req.Options)
		if err != nil {
			utils.ServerErrorWithMsg(c, "选项序列化失败")
			return
		}
		answer := strconv.Itoa(req.Answer)
		userID := c.GetInt("user_id")

		question := model.Question{
			QuestionType: string(req.QuestionType),
			Language:     req.Language,
			AiModel:      string(req.AiModel),
			Keywords:     req.Keywords,
			Title:        req.Title,
			// 序列化为json
			Options:     string(optionBytes),
			Answer:      answer,
			Explanation: req.Explanation,
			UserID:      userID,
		}
		questions = append(questions, question)
	}
	// 保存题目
	err := q.QuestionService.ConfirmQuestions(c.Request.Context(), &questions)
	if err != nil {
		utils.ServerErrorWithMsg(c, "保存题目失败")
		return
	}
	utils.Ok(c)
}

// ListQuestions 根据查询条件分页查询题目
func (q *QuestionController) ListQuestions(c *gin.Context) {
	var req dto.ListQuestionsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
	}
	// 分页的默认值处理
	page := utils.NewPage(req.PageNum, req.PageSize)
	req.PageSize = page.PageSize
	req.PageNum = page.PageNum
	// 查询题目列表(管理员查询全部的题目)
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	if role == "admin" {
		userID = 0
	}
	questions, total, err := q.QuestionService.ListQuestions(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ServerErrorWithMsg(c, "获取题目失败")
		return
	}
	// 转为res
	var list = make([]dto.QuestionRes, 0, len(questions))
	for _, question := range questions {
		var options []dto.Option
		// 将json格式的options数组转为dto
		err := json.Unmarshal([]byte(question.Options), &options)
		if err != nil {
			utils.ServerErrorWithMsg(c, "选项反序列化失败")
			return
		}
		answer, err := strconv.Atoi(question.Answer)
		if err != nil {
			utils.ServerErrorWithMsg(c, "答案转换失败")
			return
		}
		ques := dto.Question{
			Options:     options,
			Answer:      answer,
			Explanation: question.Explanation,
			Title:       question.Title,
		}
		questionRes := dto.QuestionRes{
			ID:           question.ID,
			Question:     ques,
			QuestionType: question.QuestionType,
			Language:     question.Language,
			AiModel:      question.AiModel,
			Keywords:     question.Keywords,
			CreateAt:     question.CreatedAt.Format("2006-01-02 15:04:05"),
			UserID:       question.UserID,
			UserName:     question.User.Username,
		}
		list = append(list, questionRes)
	}
	utils.SuccessMsg(c, utils.NewPageResult(list, total, req.PageNum, req.PageSize), "获取题目成功")
}

// UpdateQuestion 更新题目
func (q *QuestionController) UpdateQuestion(c *gin.Context) {
	// 获取路径参数
	questionIDStr := c.Param("question_id")
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		utils.BadRequestWithMsg(c, "无效的题目ID")
		return
	}
	// 验证权限
	permissionFlag := q.QuestionService.CheckQuestionPermission(c.Request.Context(), userID, questionID)
	// 管理员可修改所有题目
	if !permissionFlag && role != "admin" {
		utils.NotPermission(c)
		return
	}
	// 绑定参数
	var req dto.UpdateQuestionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestWithMsg(c, err.Error())
		return
	}
	// 更新题目
	err = q.QuestionService.UpdateQuestion(c.Request.Context(), userID, questionID, req)
	if err != nil {
		utils.ServerErrorWithMsg(c, "更新题目失败"+err.Error())
		return
	}
	utils.Ok(c)
}

// DeleteQuestion 删除单个题目
func (q *QuestionController) DeleteQuestion(c *gin.Context) {
	// 获取路径参数
	questionIDStr := c.Param("question_id")
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		utils.BadRequestWithMsg(c, "无效的题目ID")
		return
	}
	// 验证权限
	permissionFlag := q.QuestionService.CheckQuestionPermission(c.Request.Context(), userID, questionID)
	// 管理员角色可以删除所有题目
	if !permissionFlag && role != "admin" {
		utils.NotPermission(c)
		return
	}
	// 删除题目
	err = q.QuestionService.DeleteQuestion(c.Request.Context(), questionID)
	if err != nil {
		utils.ServerErrorWithMsg(c, "删除题目失败"+err.Error())
		return
	}
	utils.Ok(c)
}
