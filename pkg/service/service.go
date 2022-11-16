package service

import (
	"test"
	"test/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user test.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	CheckUser(username string) error
	Testing(name string) (string, error)
}

type Post interface {
	Create(userId int, post test.Post) (int, error)
	Get() ([]test.Post, error)
	GetById(id int) (test.Post, error)
	Update(userId, id int, post test.Post) error
	Delete(userId, id int) error
	GetByUserId(userId int) ([]test.Post, error)
}

type Comment interface {
	Create(postId int, comment test.Comment) (int, error)
	Get(postId int) ([]test.Comment, error)
	Update(postId, id int, comment test.Comment) error
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
		Comment: NewCommentService(repos.Comment),
	}
}
