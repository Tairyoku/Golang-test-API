package repository

import (
	"fmt"
	"gorm.io/gorm"
	"test"
)

type CommentMySql struct {
	db *gorm.DB
}

func NewCommentMySql(db *gorm.DB) *CommentMySql {
	return &CommentMySql{db: db}
}

func (p *CommentMySql) Create(postId int, comment test.Comment) (int, error) {
	tx := p.db.Begin()
	errPost := tx.Select(CommentsTable, "body", "user_id").Create(&comment).Error
	if errPost != nil {
		tx.Rollback()
		return 0, errPost
	}
	var commentsList test.CommentsList
	commentsList.PostId = postId
	commentsList.CommentId = comment.Id
	errPostsList := tx.Table(CommentsList).Select(CommentsList, "comment_id", "post_id").Create(&commentsList).Error
	if errPostsList != nil {
		tx.Rollback()
		return 0, errPostsList
	}

	return comment.Id, tx.Commit().Error
}

func (p *CommentMySql) Get(postId int) ([]test.Comment, error) {
	var comments []test.Comment
	query := fmt.Sprintf("SELECT cmt.id, cmt.user_id, u.name, cmt.body FROM %s cmt INNER JOIN %s u on u.id = cmt.user_id INNER JOIN %s cl on cl.comment_id = cmt.id WHERE cl.post_id = %d",
		CommentsTable, UsersTable, CommentsList, postId)
	err := p.db.Raw(query).Scan(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (p *CommentMySql) Update(postId, id int, comment test.Comment) error {
	var checkComment test.Comment
	query := fmt.Sprintf("SELECT * FROM %s cmt INNER JOIN %s cl on cmt.id = cl.comment_id WHERE cl.post_id = %d and cmt.id = %d",
		CommentsTable, CommentsList, postId, id)
	errCheck := p.db.Raw(query).Scan(&checkComment).Error
	if errCheck != nil {
		return errCheck
	}

	err := p.db.Select(CommentsTable, "body").Where("id = ?", id).Updates(&comment).Error
	return err
}

func (p *CommentMySql) Delete(postId, id int) error {

	var checkComment test.Comment
	query := fmt.Sprintf("SELECT * FROM %s cmt INNER JOIN %s cl on cmt.id = cl.comment_id WHERE cl.post_id = %d and cmt.id = %d",
		CommentsTable, CommentsList, postId, id)
	errCheck := p.db.Raw(query).Scan(&checkComment).Error
	if errCheck != nil {
		return errCheck
	}

	tx := p.db.Begin()
	errPost := tx.Table(CommentsTable).Where("id = ?", id).Delete(&test.Comment{}).Error
	if errPost != nil {
		tx.Rollback()
		return errPost
	}

	errPostsList := tx.Table(CommentsList).Where("comment_id = ?", id).Delete(&test.CommentsList{}).Error
	if errPostsList != nil {
		tx.Rollback()
		return errPostsList
	}
	return tx.Commit().Error
}
