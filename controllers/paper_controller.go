package controllers

import (
	"aiquiz/models/dto"
	"aiquiz/services"
	"aiquiz/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type PaperController struct {
	PaperService *services.PaperService
}

func NewPaperController(paperService *services.PaperService) *PaperController {
	return &PaperController{
		PaperService: paperService,
	}
}

// CreatePaper 创建试卷
func (p *PaperController) CreatePaper(c *gin.Context) {
	var req dto.GeneratePaperReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}
	if err := p.PaperService.GeneratePaper(c.Request.Context(), c.GetInt("user_id"), &req); err != nil {
		utils.ServerErrorWithMsg(c, "创建试卷失败"+err.Error())
		return
	}
	utils.Ok(c)
}

// ListPapers 条件分页获取试卷列表
func (p *PaperController) ListPapers(c *gin.Context) {
	var req dto.PaperListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ParamError(c)
		return
	}
	// 默认分页设置
	req.Page = utils.NewPage(req.PageNum, req.PageSize)
	// 管理员查询所有试卷
	userID := c.GetInt("user_id")
	role := c.GetString("role")
	if role == "admin" {
		userID = 0
	}
	// 查询试卷列表
	papers, total, err := p.PaperService.ListPapers(c.Request.Context(), userID, &req)
	if err != nil {
		utils.ServerErrorWithMsg(c, "获取试卷失败")
		return
	}
	// 转为res
	var list = make([]dto.PaperListRes, 0, len(papers))
	for _, paper := range papers {
		list = append(list, dto.PaperListRes{
			ID:          paper.ID,
			Title:       paper.Title,
			Description: paper.Description,
			TotalScore:  paper.TotalScore,
			CreatedAt:   paper.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   paper.UpdatedAt.Format("2006-01-02 15:04:05"),
			UserID:      paper.Creator.ID,
			UserName:    paper.Creator.Username,
		})
	}
	utils.SuccessMsg(c, utils.NewPageResult(list, total, req.PageNum, req.PageSize), "获取试卷成功")
}

// AddPaperQuestions 试卷中添加题目
func (p *PaperController) AddPaperQuestions(c *gin.Context) {
	var req []dto.AddPaperQuestionsReq
	userID := c.GetInt("user_id")
	paperID := c.GetInt("paper_id")
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}
	if len(req) == 0 {
		utils.ParamError(c)
		return
	}
	if err := p.PaperService.AddPaperQuestions(c.Request.Context(), userID, paperID, req); err != nil {
		utils.ServerErrorWithMsg(c, "添加题目失败"+err.Error())
		return
	}
	utils.Ok(c)
}

// UpdatePaper 更新试卷信息
func (p *PaperController) UpdatePaper(c *gin.Context) {
	paperID := c.GetInt("paper_id")
	var req dto.UpdatePaperReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}
	// 三者都为空，则返回参数错误
	if req.Title == "" && req.Description == "" && req.TotalScore == 0 {
		utils.ParamError(c)
		return
	}
	if req.TotalScore < 0 {
		utils.BadRequestWithMsg(c, "总分不能小于0")
		return
	}
	// 更新试卷信息
	err := p.PaperService.UpdatePaper(c.Request.Context(), paperID, req)
	if err != nil {
		utils.ServerErrorWithMsg(c, "更新试卷失败"+err.Error())
		return
	}
	utils.Ok(c)
}

// GetPaper 获取试卷详情（包含所有题目）
func (p *PaperController) GetPaper(c *gin.Context) {
	paper, err := p.PaperService.GetPaper(c.Request.Context(), c.GetInt("paper_id"))
	if err != nil {
		utils.ServerErrorWithMsg(c, "获取试卷失败"+err.Error())
		return
	}
	// 构建返回结构体
	var paperRes dto.PaperDetailRes
	// 先转换题目列表
	var questions []dto.QuestionRes
	for _, paperQues := range paper.Questions {
		// 获取具体题目
		q := paperQues.Question

		// 将json数组反序列化为option数组
		optionsStr := q.Options
		var options []dto.Option
		err := json.Unmarshal([]byte(optionsStr), &options)
		if err != nil {
			utils.ServerErrorWithMsg(c, "选项反序列化失败")
			return
		}
		answer, err := strconv.Atoi(q.Answer)
		if err != nil {
			utils.ServerErrorWithMsg(c, "答案转换失败")
			return
		}
		// 构建题目res
		ques := dto.Question{
			Options:     options,
			Answer:      answer,
			Explanation: q.Explanation,
			Title:       q.Title,
		}
		questions = append(questions, dto.QuestionRes{
			ID:           q.ID,
			QuestionType: q.QuestionType,
			Language:     q.Language,
			AiModel:      q.AiModel,
			Keywords:     q.Keywords,
			Question:     ques,
		})
	}

	paperRes = dto.PaperDetailRes{
		Questions: questions,
		PaperListRes: dto.PaperListRes{
			ID:          paper.ID,
			Title:       paper.Title,
			Description: paper.Description,
			TotalScore:  paper.TotalScore,
			CreatedAt:   paper.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   paper.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	utils.SuccessMsg(c, paperRes, "获取试卷成功")
}

// DeletePaper 删除试卷以及其题目关联
func (p *PaperController) DeletePaper(c *gin.Context) {
	paperID := c.GetInt("paper_id")
	if err := p.PaperService.DeletePaper(c.Request.Context(), paperID); err != nil {
		utils.ServerErrorWithMsg(c, "删除试卷失败"+err.Error())
		return
	}
	utils.Ok(c)
}

// DeletePaperQuestions 删除试卷中的题目
func (p *PaperController) DeletePaperQuestions(c *gin.Context) {
	questionIDStr := c.Param("question_id")
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		utils.BadRequestWithMsg(c, "无效的题目ID")
	}
	// 根据paperID和questionID删除题目
	deleted, err := p.PaperService.DeletePaperQuestion(c.Request.Context(), c.GetInt("paper_id"), questionID)
	if err != nil {
		utils.ServerErrorWithMsg(c, "删除题目失败")
		return
	}
	// 删除条目必须大于0
	if !deleted {
		utils.FailMsg(c, utils.ERROR_RECORD_NOT_EXIST, "没有找到要删除的记录")
		return
	}
	utils.Ok(c)
}

func (p *PaperController) UpdatePaperQuestionOrder(c *gin.Context) {
	var req []dto.QuestionOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c)
		return
	}
	// 参数不能为空
	if len(req) == 0 {
		utils.BadRequestWithMsg(c, "参数列表不能为空")
		return
	}
	// 校验新order在列表内唯一
	orderSet := make(map[int]struct{})
	for _, q := range req {
		if _, exists := orderSet[q.QuestionOrder]; exists {
			utils.BadRequestWithMsg(c, fmt.Sprintf("存在重复的order值：%d", q.QuestionOrder))
			return
		}
		orderSet[q.QuestionOrder] = struct{}{}
	}
	// 更新题目顺序
	if err := p.PaperService.UpdatePaperQuestionOrder(c.Request.Context(), c.GetInt("paper_id"), req); err != nil {
		utils.ServerErrorWithMsg(c, "更新题目顺序失败"+err.Error())
		return
	}
	utils.Ok(c)
}
