package dao

import (
	"aiquiz/dao/model"
	"aiquiz/utils"
	"context"
	"gorm.io/gorm"
)

type QuestionDao struct {
	DB *gorm.DB
}

func NewQuestionDAO(db *gorm.DB) *QuestionDao {
	return &QuestionDao{
		DB: db,
	}
}

func (dao *QuestionDao) AddQuestions(c context.Context, questions *[]model.Question) error {
	return dao.DB.WithContext(c).CreateInBatches(questions, len(*questions)).Error
}

func (dao *QuestionDao) ListQuestions(
	c context.Context,
	userID int,
	title, questionType, keywords, language, aiModel string,
	page utils.Page,
) ([]model.Question, int64, error) {

	var questions []model.Question
	// 构建查询条件
	query := dao.DB.WithContext(c).Model(&model.Question{}).Order("created_at desc")
	if userID != 0 {
		query = query.Where("user_id = ?", userID)
	}
	if title != "" {
		query = query.Where("title LIKE?", "%"+title+"%")
	}
	if questionType != "" {
		query = query.Where("question_type =?", questionType)
	}
	if keywords != "" {
		query = query.Where("keywords LIKE?", "%"+keywords+"%")
	}
	if language != "" {
		query = query.Where("language =?", language)
	}
	if aiModel != "" {
		query = query.Where("ai_model =?", aiModel)
	}
	// 查询总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 使用分页器
	err = query.Scopes(utils.Paginate(page)).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Find(&questions).Error

	if err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

func (dao *QuestionDao) UpdateQuestion(c context.Context, q *model.Question) error {
	return dao.DB.WithContext(c).Model(&model.Question{}).Where("id = ?", q.ID).Updates(q).Error
}

func (dao *QuestionDao) QueryQuestion(c context.Context, userID, questionID int) (*model.Question, error) {
	var question model.Question
	err := dao.DB.WithContext(c).Model(&model.Question{}).Where("id = ?", questionID).Where("user_id = ?", userID).Take(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (dao *QuestionDao) DeleteQuestion(c context.Context, questionID int) error {
	return dao.DB.WithContext(c).Model(&model.Question{}).Delete(&model.Question{
		ID: questionID,
	}).Error
}

// GetExistingQuestionIDs 返回存在于数据库中的题目 ID 列表
func (dao *QuestionDao) GetExistingQuestionIDs(c context.Context, userID int, questionIDList []int) ([]int, error) {
	var existingIDs []int
	err := dao.DB.WithContext(c).Model(&model.Question{}).Select("id").Where("user_id = ?", userID).
		Where("id IN ?", questionIDList).Find(&existingIDs).Error
	if err != nil {
		return nil, err
	}
	return existingIDs, nil

}

func (dao *QuestionDao) DeleteQuestionByUserID(c context.Context, tx *gorm.DB, userID int) error {
	return tx.WithContext(c).
		Where("user_id = ?", userID).
		Delete(&model.Question{}).Error
}
