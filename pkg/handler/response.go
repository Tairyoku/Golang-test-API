package handler

import (
	"github.com/labstack/echo/v4"
	"test/pkg/repository/models"
)

type GetPostsResponse struct {
	Posts []models.Post `json:"posts"`
}

type GetCommentsResponse struct {
	Comments []models.Comment `json:"comments"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type IdResponse struct {
	Id int `json:"id"`
}

type PostRequest struct {
	Title string `json:"title" form:"title" binding:"required"`
	Anons string `json:"anons" form:"anons" binding:"required"`
}

type CommentRequest struct {
	Body string `json:"body"  binding:"required"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type UserResponse struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Username string `json:"username" form:"username"  binding:"required"`
	Password string `json:"password" gorm:"column:password_hash" form:"password"  binding:"required"`
}

func NewErrorResponse(c echo.Context, statusCode int, message string) {
	errRes := c.JSON(statusCode, ErrorResponse{Message: message})
	if errRes != nil {
		return
	}
}
