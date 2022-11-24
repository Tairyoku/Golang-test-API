package service

import (
	"test/pkg/repository"
	"test/pkg/repository/models"
)

type PostService struct {
	repository repository.Post
}

func NewPostService(repository repository.Post) *PostService {
	return &PostService{repository: repository}
}

func (p *PostService) Create(post models.Post) (int, error) {
	return p.repository.Create(post)
}

func (p *PostService) Get() ([]models.Post, error) {
	return p.repository.Get()
}

func (p *PostService) GetById(id int) (models.Post, error) {
	return p.repository.GetById(id)
}

func (p *PostService) GetByUserId(userId int) ([]models.Post, error) {
	return p.repository.GetByUserId(userId)
}

func (p *PostService) Update(id int, post models.Post) error {
	return p.repository.Update(id, post)
}

func (p *PostService) Delete(id int) error {
	return p.repository.Delete(id)
}
