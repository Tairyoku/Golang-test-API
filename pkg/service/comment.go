package service

import (
	"test/pkg/repository"
	"test/pkg/repository/models"
)

type CommentService struct {
	repository repository.Comment
}

func NewCommentService(repository repository.Comment) *CommentService {
	return &CommentService{repository: repository}
}

func (p *CommentService) Create(comment models.Comment) (int, error) {
	return p.repository.Create(comment)
}

func (p *CommentService) Get(postId int) ([]models.Comment, error) {
	return p.repository.Get(postId)
}

func (p *CommentService) Update(postId, id int, comment models.Comment) error {
	return p.repository.Update(postId, id, comment)
}

func (p *CommentService) Delete(postId, id int) error {
	return p.repository.Delete(postId, id)
}
