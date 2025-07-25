package services

import (
	"aiquiz/dao"
	"aiquiz/dao/model"
	"aiquiz/models/dto"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type UserService struct {
	userDao     *dao.UserDao
	questionDao *dao.QuestionDao
	paperDao    *dao.PaperDao
}

func NewUserService(userDAO *dao.UserDao, questionDao *dao.QuestionDao, paperDao *dao.PaperDao) *UserService {
	return &UserService{
		userDao:     userDAO,
		questionDao: questionDao,
		paperDao:    paperDao,
	}
}
func (s *UserService) Create(c context.Context, user *model.User) error {
	return s.userDao.Create(c, user)
}

func (s *UserService) Login(c context.Context, user *model.User) (*model.User, error) {
	return s.userDao.ValidateLogin(c, user.Username, user.PasswordHash)
}

func (s *UserService) ListUsersByPage(c context.Context, req *dto.UserListReq) ([]model.User, int64, error) {
	return s.userDao.ListUsers(c, req.Username, req.Page)
}

func (s *UserService) UpdateUser(ctx context.Context, id int, username string, password string) error {
	return s.userDao.UpdateUser(ctx, id, username, password)
}

func (s *UserService) DeleteUser(c context.Context, deletedUserID int) error {
	return s.userDao.DB.Transaction(func(tx *gorm.DB) error {
		// 删除试卷题目关联表以及试卷表数据
		err := s.paperDao.DeletePaperByUserID(c, tx, deletedUserID)
		if err != nil {
			return fmt.Errorf("删除试卷失败: %w", err)
		}
		// 删除问题表数据
		err = s.questionDao.DeleteQuestionByUserID(c, tx, deletedUserID)
		if err != nil {
			return fmt.Errorf("删除题目失败: %w", err)
		}
		// 删除用户表数据
		err = s.userDao.DeleteUser(c, tx, deletedUserID)
		if err != nil {
			return fmt.Errorf("删除用户失败: %w", err)
		}
		return nil
	})
}
