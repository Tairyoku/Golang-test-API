package repository

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	UsersTable    = "users"
	PostsTable    = "posts"
	CommentsTable = "comments"
	CommentsList  = "comments_list"
	PostsList     = "posts_list"
	CommentsUser  = "comments_user_list"
)

type Config struct {
	Username string
	Password string
	Host     string
	Url      string
	DBName   string
}

func NewMysqlDB(cnf Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s@%s%s(%s)/%s", cnf.Username, cnf.Password, cnf.Host, cnf.Url, cnf.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
