package models

type Post struct {
	Id     int    `json:"id" gorm:"<-:false"`
	UserId int    `json:"user_id"`
	Title  string `json:"title" form:"title" binding:"required"`
	Anons  string `json:"anons" form:"anons" binding:"required"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}
