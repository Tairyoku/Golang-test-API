package models

type Comment struct {
	Id     int    `json:"id"  gorm:"<-:false"`
	PostId int    `json:"post_id"`
	UserId int    `json:"user_id"`
	Body   string `json:"body"  binding:"required"`
}
