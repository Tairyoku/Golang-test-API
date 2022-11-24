package repository

import (
	"fmt"
	"gorm.io/gorm"
	"test/pkg/repository/models"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (p *CommentRepository) Create(comment models.Comment) (int, error) {
	tx := p.db.Begin()
	errPost := tx.Select(CommentsTable, "body", "user_id", "post_id").Create(&comment).Error
	return comment.Id, errPost
}

func (p *CommentRepository) Get(postId int) ([]models.Comment, error) {
	var comments []models.Comment
	query := fmt.Sprintf("SELECT * FROM %s cmt WHERE cmt.post_id = %d",
		CommentsTable, postId)
	err := p.db.Raw(query).Scan(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (p *CommentRepository) Update(postId, id int, comment models.Comment) error {
	err := p.db.Select(CommentsTable, "body").Where("id = ? and post_id = ?", id, postId).Updates(&comment).Error
	return err
}

func (p *CommentRepository) Delete(postId, id int) error {

	errPost := p.db.Table(CommentsTable).Where("id = ? and post_id = ?", id, postId).Delete(&models.Comment{}).Error
	return errPost

}
