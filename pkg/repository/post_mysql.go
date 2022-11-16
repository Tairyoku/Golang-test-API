package repository

import (
	"fmt"
	"gorm.io/gorm"
	"test"
)

type PostMySql struct {
	db *gorm.DB
}

func NewPostMySql(db *gorm.DB) *PostMySql {
	return &PostMySql{db: db}
}

func (p *PostMySql) GetByUserId(userId int) ([]test.Post, error) {
	var posts []test.Post
	query := fmt.Sprintf("SELECT * FROM %s post INNER JOIN %s ul on post.id = ul.post_id WHERE ul.user_id = %d",
		PostsTable, PostsList, userId)
	err := p.db.Raw(query).Scan(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *PostMySql) Create(userId int, post test.Post) (int, error) {
	tx := p.db.Begin()
	errPost := tx.Select(PostsTable, "title", "anons").Create(&post).Error
	if errPost != nil {
		tx.Rollback()
		return 0, errPost
	}
	var postsList test.PostsList
	postsList.PostId = post.Id
	postsList.UserId = userId
	errPostsList := tx.Table(PostsList).Select(PostsList, "user_id", "post_id").Create(&postsList).Error
	if errPostsList != nil {
		tx.Rollback()
		return 0, errPostsList
	}
	return post.Id, tx.Commit().Error
}

func (p *PostMySql) Get() ([]test.Post, error) {
	var posts []test.Post
	err := p.db.Find(&posts).Error
	return posts, err
}

func (p *PostMySql) GetById(id int) (test.Post, error) {
	var post test.Post
	err := p.db.Table(PostsTable).First(&post, id).Error
	return post, err
}

func (p *PostMySql) Update(userId, id int, post test.Post) error {
	var checkPost test.Post
	query := fmt.Sprintf("SELECT * FROM %s post INNER JOIN %s ul on post.id = ul.post_id WHERE ul.user_id = %d and ul.post_id = %d",
		PostsTable, PostsList, userId, id)
	errCheck := p.db.Raw(query).Scan(&checkPost).Error
	if errCheck != nil {
		return errCheck
	}

	err := p.db.Select(PostsTable, "title", "anons").Where("id = ?", id).Updates(post).Error
	return err
}

func (p *PostMySql) Delete(userId, id int) error {

	var checkPost test.Post
	query := fmt.Sprintf("SELECT * FROM %s post INNER JOIN %s ul on post.id = ul.post_id WHERE ul.user_id = %d and ul.post_id = %d",
		PostsTable, PostsList, userId, id)
	errCheck := p.db.Raw(query).Scan(&checkPost).Error
	if errCheck != nil {
		return errCheck
	}

	tx := p.db.Begin()
	errPost := tx.Table(PostsTable).Where("id = ?", id).Delete(&test.Post{}).Error
	if errPost != nil {
		tx.Rollback()
		return errPost
	}

	errPostsList := tx.Table(PostsList).Where("post_id = ?", id).Delete(&test.PostsList{}).Error
	if errPostsList != nil {
		tx.Rollback()
		return errPostsList
	}
	return tx.Commit().Error
}
