package services

import (
	"aiquiz/dao"
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"context"
	"errors"
)

type PaperService struct {
	paperDao    *dao.PaperDao
	questionDao *dao.QuestionDao
}

func NewPaperService(paperDao *dao.PaperDao, questionDao *dao.QuestionDao) *PaperService {
	return &PaperService{
		paperDao:    paperDao,
		questionDao: questionDao,
	}
}
func (s *PaperService) GeneratePaper(c context.Context, userID int, req *dto.GeneratePaperReq) error {
	// 转为model
	paper := &model.Paper{
		Title:       req.Title,
		Description: req.Description,
		TotalScore:  req.TotalScore,
		CreatorID:   userID,
	}
	// 创建试卷
	return s.paperDao.GeneratePaper(c, paper)
}

func (s *PaperService) ListPapers(c context.Context, userID int, req *dto.PaperListReq) ([]model.Paper, int64, error) {
	return s.paperDao.ListPapers(c, userID, req.Title, req.Description, req.Page)
}

func (s *PaperService) AddPaperQuestions(c context.Context, userID, paperID int, req []dto.AddPaperQuestionsReq) error {
	questionIDList := make([]int, 0, len(req))
	for _, question := range req {
		questionIDList = append(questionIDList, question.QuestionID)
	}
	// 验证所有题目ID是否存在
	existQuestionIds, err := s.questionDao.GetExistingQuestionIDs(c, userID, questionIDList)
	if err != nil {
		return err
	}
	if len(existQuestionIds) != len(questionIDList) {
		return errors.New("某些题目不存在")
	}
	// 将req转换为model
	var paperQuestions = make([]model.PaperQuestion, 0, len(req))
	for _, question := range req {
		paperQuestions = append(paperQuestions, model.PaperQuestion{
			PaperID:    paperID,
			QuestionID: question.QuestionID,
			Score:      question.Score,
		})
	}

	// 插入试卷题目(需要获取目前最大order再插入)
	err = s.paperDao.AddPaperQuestions(c, paperID, paperQuestions)
	if err != nil {
		return err
	}
	return nil
}

func (s *PaperService) UpdatePaper(c context.Context, paperID int, req dto.UpdatePaperReq) error {
	return s.paperDao.UpdatePaper(c, paperID, req.Title, req.Description, req.TotalScore)
}

func (s *PaperService) DeletePaperQuestion(c context.Context, paperID, questionID int) (bool, error) {
	return s.paperDao.DeletePaperQuestion(c, paperID, questionID)
}
func (s *PaperService) GetPaper(c context.Context, paperID int) (*model.Paper, error) {
	return s.paperDao.GetPaper(c, paperID)
}

func (s *PaperService) DeletePaper(c context.Context, paperID int) error {
	return s.paperDao.DeletePaper(c, paperID)
}

func (s *PaperService) UpdatePaperQuestionOrder(c context.Context, paperID int, req []dto.QuestionOrderReq) error {
	// 查询试卷全部题目
	existQuestionIDs, err := s.paperDao.GetPaperQuestionIDs(paperID)
	if err != nil {
		return err
	}
	// 判断传入的题目ID是否存在
	if len(req) != len(existQuestionIDs) {
		return errors.New("试卷题目数量与参数不一致")
	}
	// 将已存在的题目ID的List转换为map
	existIDMap := make(map[int]bool)
	for _, id := range existQuestionIDs {
		existIDMap[id] = true
	}

	// 检查req中的每个题目ID是否都存在于existQuestionIDs中
	var paperQuestions = make([]model.PaperQuestion, 0, len(req))
	for _, question := range req {
		paperQuestions = append(paperQuestions, model.PaperQuestion{
			PaperID:       paperID,
			QuestionID:    question.QuestionID,
			QuestionOrder: question.QuestionOrder,
			Score:         question.Score,
		})
		if !existIDMap[question.QuestionID] {
			return errors.New("某些题目不存在")
		}
	}
	// 执行批量更新
	return s.paperDao.UpdatePaperQuestionOrder(c, paperID, paperQuestions)

}
