package repository

import (
	"gorm.io/gorm"
	"test"
)

type Authorization interface {
	CreateUser(user test.User) (int, error)
	GetUser(username, password string) (test.User, error)
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

type Repository struct {
	Authorization
	Post
	Comment
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthMySql(db),
		Post:          NewPostMySql(db),
		Comment:       NewCommentMySql(db),
	}
}
