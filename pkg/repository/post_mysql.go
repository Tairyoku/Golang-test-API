package repository

import (
	"fmt"
	"gorm.io/gorm"
	"test/pkg/repository/models"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}
func (p *PostRepository) Get() ([]models.Post, error) {
	var posts []models.Post
	err := p.db.Find(&posts).Error
	return posts, err
}

func (p *PostRepository) GetById(id int) (models.Post, error) {
	var post models.Post
	err := p.db.Table(PostsTable).First(&post, id).Error
	return post, err
}

func (p *PostRepository) GetByUserId(userId int) ([]models.Post, error) {
	var posts []models.Post
	query := fmt.Sprintf("SELECT * FROM %s post WHERE post.user_id = %d",
		PostsTable, userId)
	err := p.db.Raw(query).Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostRepository) Create(post models.Post) (int, error) {
	errPost := p.db.Select(PostsTable, "user_id", "title", "anons").Create(&post).Error
	return post.Id, errPost
}

func (p *PostRepository) Update(id int, post models.Post) error {
	err := p.db.Select(PostsTable, "title", "anons").Where("id = ?", id).Updates(&post).Error
	return err
}

func (p *PostRepository) Delete(id int) error {
	errPost := p.db.Table(PostsTable).Where("id = ?", id).Delete(&models.Post{}).Error
	return errPost
}
