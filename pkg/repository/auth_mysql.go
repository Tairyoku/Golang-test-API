package repository

import (
	"gorm.io/gorm"
	"test/pkg/repository/models"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) CreateUser(user models.User) (int, error) {
	err := a.db.Select(UsersTable, "name", "username", "password_hash").Create(&user).Error
	if user.Id == 0 {
		return 0, err
	}
	return user.Id, nil
}

func (a *AuthRepository) Testing(name string) (string, error) {
	if name != "name" {
		return "error", nil
	}
	return name, nil
}

func (a *AuthRepository) GetUser(username, password string) (models.User, error) {
	var user models.User
	err := a.db.Where("username = ? and password_hash = ?", username, password).Find(&user).Error
	return user, err
}

func (a *AuthRepository) CheckUser(username string) error {
	var user models.User
	err := a.db.Table(UsersTable).Where("username = ?", username).First(&user).Error
	return err
}
