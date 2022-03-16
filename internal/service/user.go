package service

import (
	"errors"

	"github.com/sendya/pkg/encode/passwd"
	"gorm.io/gorm"

	"example/internal/models"
)

var (
	ErrUserNotFound      = errors.New("user_not_found")
	ErrUserHasDisabled   = errors.New("user_has_disabled")
	ErrUserPasswordWrong = errors.New("user_password_wrong")
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db,
	}
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	var account models.User
	if err := s.db.Model(&models.User{}).
		Where(&models.User{
			Username: username,
		}).Find(&account).Error; err != nil {
		return nil, ErrUserNotFound
	}

	// 密码加盐进行校验
	if passwd.Sum(password, account.Salt) != account.Password {
		return nil, ErrUserPasswordWrong
	}

	return &account, nil
}
