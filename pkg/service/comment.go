package service

import (
	"test"
	"test/pkg/repository"
)

type CommentService struct {
	repository repository.Comment
}

func NewCommentService(repository repository.Comment) *CommentService {
	return &CommentService{repository: repository}
}



func (p *CommentService) Create(postId int, comment test.Comment) (int, error) {
	return p.repository.Create(postId, comment)
}

func (p *CommentService) Get(postId int) ([]test.Comment, error) {
	return p.repository.Get(postId)
}

func (p *CommentService) Update(postId, id int, comment test.Comment) error {
	return p.repository.Update(postId, id, comment)
}

func (p *CommentService) Delete(postId, id int) error {
	return p.repository.Delete(postId, id)
}
