package test

type Post struct {
	Id    int    `json:"id" gorm:"<-:false"`
	Title string `json:"title" form:"title" binding:"required"`
	Anons string `json:"anons" form:"anons" binding:"required"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type PostsList struct {
	Id     int
	UserId int `json:"user_id"`
	PostId int `json:"post_id"`
}
