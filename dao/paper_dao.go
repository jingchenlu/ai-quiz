package dao

import (
	"aiquiz/dao/model"
	"aiquiz/utils"
	"context"
	"errors"
	"gorm.io/gorm"
)

type PaperDao struct {
	DB *gorm.DB
}

func NewPaperDAO(db *gorm.DB) *PaperDao {
	return &PaperDao{DB: db}
}

func (dao *PaperDao) GeneratePaper(c context.Context, paper *model.Paper) error {
	return dao.DB.WithContext(c).Create(paper).Error
}

func (dao *PaperDao) GetPaper(c context.Context, paperID int) (*model.Paper, error) {
	var paper model.Paper
	err := dao.DB.WithContext(c).Model(&model.Paper{}).
		Where("id = ?", paperID).
		// 预加载paperQuestion一对多的关联
		Preload("Questions").
		// 预加载题目一对一关联
		Preload("Questions.Question").
		First(&paper).Error
	return &paper, err
}

func (dao *PaperDao) ListPapers(c context.Context, userID int, title, description string, page utils.Page) (papers []model.Paper, total int64, err error) {
	query := dao.DB.WithContext(c).Model(&model.Paper{}).Order("created_at desc")
	if userID != 0 {
		query = query.Where("creator_id = ?", userID)
	}
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if description != "" {
		query = query.Where("description LIKE ?", "%"+description+"%")
	}
	// 查询总数
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 使用分页器
	err = query.Scopes(utils.Paginate(page)).Preload("Creator", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Find(&papers).Error

	if err != nil {
		return nil, 0, err
	}
	return papers, total, nil
}

// AddPaperQuestions 试卷中添加题目
func (dao *PaperDao) AddPaperQuestions(c context.Context, paperID int, paperQuestions []model.PaperQuestion) error {
	// 开启事务
	return dao.DB.Transaction(func(tx *gorm.DB) error {
		var maxOrder *int
		// 直接获取最大题目顺序，通过乐观锁来插入（利用db唯一索引）（高并发场景 可以使用mysql的for update行锁获取最大顺序）
		err := tx.WithContext(c).Raw("SELECT MAX(`question_order`) FROM paper_questions WHERE paper_id = ?", paperID).Scan(&maxOrder).Error
		if err != nil {
			return err
		}
		// 说明试卷目前没有题目，初始化maxOrder为0
		if maxOrder == nil {
			maxOrder = new(int)
			*maxOrder = 0
		}
		// 从 maxOrder+1 开始插入
		for i := range paperQuestions {
			paperQuestions[i].QuestionOrder = *maxOrder + i + 1
		}
		return tx.WithContext(c).Model(&model.PaperQuestion{}).Create(&paperQuestions).Error
	})
}

func (dao *PaperDao) UpdatePaper(c context.Context, paperID int, title string, description string, score int) error {
	return dao.DB.WithContext(c).Model(&model.Paper{}).Where("id = ?", paperID).Updates(model.Paper{
		Title:       title,
		Description: description,
		TotalScore:  score,
	}).Error
}

func (dao *PaperDao) DeletePaperQuestion(c context.Context, paperID, questionID int) (bool, error) {
	result := dao.DB.WithContext(c).
		Where("paper_id = ? AND question_id = ?", paperID, questionID).
		Delete(&model.PaperQuestion{})
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (dao *PaperDao) DeletePaper(c context.Context, paperID int) error {
	return dao.DB.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 先删除paper_questions表中的关联
		err := tx.WithContext(c).Where("paper_id = ?", paperID).Delete(&model.PaperQuestion{}).Error
		if err != nil {
			return err
		}
		// 再删除paper表
		return tx.WithContext(c).Where("id = ?", paperID).Delete(&model.Paper{}).Error
	})
}

func (dao *PaperDao) GetPaperQuestionIDs(paperID int) ([]int, error) {
	var questionIDs []int
	err := dao.DB.Model(&model.PaperQuestion{}).Select("question_id").Where("paper_id = ?", paperID).Find(&questionIDs).Error
	return questionIDs, err
}

// UpdatePaperQuestionOrder 由于有唯一索引（paperID, order），无法使用插入冲突时更新。故选用先删除再新增的方式
func (dao *PaperDao) UpdatePaperQuestionOrder(c context.Context, paperID int, questions []model.PaperQuestion) error {
	return dao.DB.WithContext(c).Transaction(func(tx *gorm.DB) error {
		// 删除该试卷下所有有效题目关联
		if err := tx.Where("paper_id = ?", paperID).Delete(&model.PaperQuestion{}).Error; err != nil {
			return errors.New("删除旧题目关联失败")
		}

		// 批量插入新的题目顺序
		if err := tx.Create(&questions).Error; err != nil {
			return errors.New("插入新题目顺序失败")
		}
		return nil
	})
}

func (dao *PaperDao) DeletePaperByUserID(c context.Context, tx *gorm.DB, userID int) error {
	// 先删除试卷题目关联表（通过试卷ID关联）
	// 查询该用户的所有试卷ID
	var paperIDs []int
	if err := tx.WithContext(c).
		Model(&model.Paper{}).
		Select("id").
		Where("creator_id = ?", userID).
		Find(&paperIDs).Error; err != nil {
		return err
	}
	// 若有试卷，批量删除关联表数据
	if len(paperIDs) > 0 {
		if err := tx.WithContext(c).
			Where("paper_id IN (?)", paperIDs).
			Delete(&model.PaperQuestion{}).Error; err != nil {
			return err
		}
	}

	// 删除试卷表数据
	return tx.WithContext(c).
		Where("creator_id = ?", userID).
		Delete(&model.Paper{}).Error
}
