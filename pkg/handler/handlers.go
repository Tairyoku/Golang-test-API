package handler

import (
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"test/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *echo.Echo {
	router := echo.New()
	router.GET("/swagger/server/*", echoSwagger.WrapHandler)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	google := router.Group("/google")
	google.GET("/login", h.GoogleLogin)
	google.GET("/callback", h.GoogleCallback)

	api := router.Group("/api")
	post := api.Group("/posts")
	{
		post.GET("", h.GetPosts)
		post.GET("/user/:id", h.GetUserPosts)
		post.GET("/:id", h.GetPostById)
		post.POST("", h.PostPost, h.userIdentify)
		post.PUT("/:id", h.UpdatePost, h.userIdentify)
		post.DELETE("/:id", h.DeletePost, h.userIdentify)
	}

	comment := post.Group("/:postId/comments", h.userIdentify)
	{
		comment.GET("", h.GetComments)
		comment.POST("", h.PostComment)
		comment.PUT("/:id", h.UpdateComment)
		comment.DELETE("/:id", h.DeleteComment)
	}
	return router
}
