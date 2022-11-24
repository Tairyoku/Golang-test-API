package repository

import (
	"gorm.io/gorm"
	"test/pkg/repository/models"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
	CheckUser(username string) error
	Testing(name string) (string, error)
}

type Post interface {
	Create(post models.Post) (int, error)
	Get() ([]models.Post, error)
	GetById(id int) (models.Post, error)
	GetByUserId(userId int) ([]models.Post, error)
	Update(id int, post models.Post) error
	Delete(id int) error
}

type Comment interface {
	Create(comment models.Comment) (int, error)
	Get(postId int) ([]models.Comment, error)
	Update(postId, id int, comment models.Comment) error
	Delete(postId, id int) error
}

type Repository struct {
	Authorization
	Post
	Comment
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		Post:          NewPostRepository(db),
		Comment:       NewCommentRepository(db),
	}
}
