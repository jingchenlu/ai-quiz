package services

import (
	"aiquiz/dao"
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"context"
	"encoding/json"
	"errors"
	"strconv"
)

type QuestionService struct {
	questionDao *dao.QuestionDao
}

func NewQuestionService(questionDAO *dao.QuestionDao) *QuestionService {
	return &QuestionService{
		questionDao: questionDAO,
	}
}

func (s *QuestionService) ConfirmQuestions(c context.Context, questions *[]model.Question) error {
	return s.questionDao.AddQuestions(c, questions)
}

func (s *QuestionService) ListQuestions(c context.Context, userID int, req *dto.ListQuestionsReq) ([]model.Question, int64, error) {
	return s.questionDao.ListQuestions(c, userID, req.Title, string(req.QuestionType), req.Keywords, req.Language, string(req.AiModel), req.Page)
}

func (s *QuestionService) UpdateQuestion(c context.Context, useID, questionID int, req dto.UpdateQuestionReq) error {
	// 构建Question
	options, err := json.Marshal(req.Options)
	if err != nil {
		return errors.New("选项序列化失败")
	}
	question := model.Question{
		ID:           questionID,
		Title:        req.Title,
		Keywords:     req.Keywords,
		Language:     req.Language,
		QuestionType: string(req.QuestionType),
		Answer:       strconv.Itoa(req.Answer),
		Explanation:  req.Explanation,
		Options:      string(options),
		UserID:       useID,
	}
	return s.questionDao.UpdateQuestion(c, &question)
}

func (s *QuestionService) CheckQuestionPermission(c context.Context, userID, questionID int) bool {
	q, err := s.questionDao.QueryQuestion(c, userID, questionID)
	return err == nil && q != nil
}

func (s *QuestionService) DeleteQuestion(c context.Context, questionID int) error {
	return s.questionDao.DeleteQuestion(c, questionID)
}
