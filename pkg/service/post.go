package service

import (
	"test"
	"test/pkg/repository"
)

type PostService struct {
	repository repository.Post
}

func NewPostService(repository repository.Post) *PostService {
	return &PostService{repository: repository}
}

func (p *PostService) Create(userId int, post test.Post) (int, error) {
	return p.repository.Create(userId, post)
}

func (p *PostService) Get() ([]test.Post, error) {
	return p.repository.Get()
}

func (p *PostService) GetById(id int) (test.Post, error) {
	return p.repository.GetById(id)
}

func (p *PostService) GetByUserId(userId int) ([]test.Post, error) {
	return p.repository.GetByUserId(userId)
}

func (p *PostService) Update(userId, id int, post test.Post) error {
	return p.repository.Update(userId, id, post)
}

func (p *PostService) Delete(userId, id int) error {
	return p.repository.Delete(userId, id)
}
