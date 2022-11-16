package repository

import (
	"gorm.io/gorm"
	"test"
)

type AuthMySql struct {
	db *gorm.DB
}

func NewAuthMySql(db *gorm.DB) *AuthMySql {
	return &AuthMySql{db: db}
}

func (a *AuthMySql) CreateUser(user test.User) (int, error) {
	err := a.db.Select(UsersTable, "name", "username", "password_hash").Create(&user).Error
	if user.Id == 0 {
		return 0, err
	}
	return user.Id, nil
}

func (a *AuthMySql) Testing(name string) (string, error) {
	if name != "name" {
		return "error", nil
	}
	return name, nil
}

func (a *AuthMySql) GetUser(username, password string) (test.User, error) {
	var user test.User
	err := a.db.Where("username = ? and password_hash = ?", username, password).Find(&user).Error
	return user, err
}

func (a *AuthMySql) CheckUser(username string) error {
	var user test.User
	err := a.db.Table(UsersTable).Where("username = ?", username).First(&user).Error
	return err
}
