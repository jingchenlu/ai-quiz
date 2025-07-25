package dao

import (
	"aiquiz/dao/model"
	"aiquiz/utils"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserDao struct {
	DB *gorm.DB
}

// NewUserDAO 创建用户DAO实例
func NewUserDAO(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// Create 创建用户
func (dao *UserDao) Create(c context.Context, user *model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	return dao.DB.WithContext(c).Create(user).Error
}

// ValidateLogin 验证用户名和密码
func (dao *UserDao) ValidateLogin(c context.Context, username, password string) (*model.User, error) {
	var user model.User
	err := dao.DB.WithContext(c).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	// 比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}
	return &user, nil
}

// ListUsers 分页查询用户
func (dao *UserDao) ListUsers(c context.Context, username string, page utils.Page) ([]model.User, int64, error) {
	var users []model.User
	query := dao.DB.WithContext(c).
		Model(&model.User{}).
		Select("id", "username", "role", "created_at", "updated_at")
	// 条件查询
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	// 按创建时间降序
	query = query.Order("created_at desc")
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 使用分页器
	query.Scopes(utils.Paginate(page))
	err = query.Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (dao *UserDao) UpdateUser(c context.Context, id int, username string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordHash := string(hashedPassword)
	return dao.DB.WithContext(c).Updates(model.User{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
	}).Error
}

func (dao *UserDao) GetUserByID(c context.Context, userID int) (*model.User, error) {
	var user model.User
	err := dao.DB.WithContext(c).Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (dao *UserDao) DeleteUser(c context.Context, tx *gorm.DB, userID int) error {
	return tx.WithContext(c).Where("id = ?", userID).Delete(&model.User{}).Error
}
