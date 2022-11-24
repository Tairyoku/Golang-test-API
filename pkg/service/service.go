package service

import (
	"test/pkg/repository"
	"test/pkg/repository/models"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	CheckUser(username string) error
	Testing(name string) (string, error)
}

type Post interface {
	Create(post models.Post) (int, error)
	Get() ([]models.Post, error)
	GetById(id int) (models.Post, error)
	Update(id int, post models.Post) error
	Delete(id int) error
	GetByUserId(userId int) ([]models.Post, error)
}

type Comment interface {
	Create(comment models.Comment) (int, error)
	Get(postId int) ([]models.Comment, error)
	Update(postId, id int, comment models.Comment) error
	Delete(postId, id int) error
}

type Service struct {
	Authorization
	Post
	Comment
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Post:          NewPostService(repos.Post),
		Comment:       NewCommentService(repos.Comment),
	}
}
