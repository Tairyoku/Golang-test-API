package test

type Comment struct {
	Id     int    `json:"id"  gorm:"<-:false"`
	UserId int    `json:"user_id" binding:"required"`
	Body   string `json:"body"  binding:"required"`
}

type CommentsList struct {
	Id        int
	PostId    int
	CommentId int
}
