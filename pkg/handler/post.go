package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"test/pkg/repository/models"
)

// GetPosts godoc
// @Summary     Find all posts
// @Description Get all posts
// @Tags        posts
// @Produce     json
// @Success     200 {object} GetPostsResponse
// @Failure 	500 {object} ErrorResponse	 "something went wrong"
// @Router      /api/posts [get]
func (h *Handler) GetPosts(c echo.Context) error {

	posts, err := h.services.Post.Get()
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "something went wrong")
		return nil
	}
	_, errEnCd := json.Marshal(&posts)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, GetPostsResponse{Posts: posts})
	if errRes != nil {
		return nil
	}
	return nil
}

// GetUserPosts godoc
// @Summary     Find all user's posts by user ID
// @Description Get user's posts by ID
// @Tags        posts
// @Produce     json
// @Param       id  path     int true "User ID"
// @Success     200 {object} GetPostsResponse
// @Failure 	400 {object} ErrorResponse	 "ID is not integer"
// @Failure 	500 {object} ErrorResponse	 "wrong user ID"
// @Router      /api/posts/user/{id} [get]
func (h *Handler) GetUserPosts(c echo.Context) error {
	userId, errParams := GetParam(c, ParamId)
	if errParams != nil {
		return nil
	}

	posts, err := h.services.Post.GetByUserId(userId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "wrong user ID")
		return nil
	}
	_, errEnCd := json.Marshal(posts)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, GetPostsResponse{
		Posts: posts,
	})
	if errRes != nil {
		return nil
	}
	return nil
}

// GetPostById godoc
// @Summary     Find post by post ID
// @Description Get post by post ID
// @Tags        posts
// @Produce     json
// @Param       id  path     int true "Post ID"
// @Success     200 {object} test.Post
// @Failure 	400 {object} ErrorResponse	 "ID is not integer"
// @Failure 	500 {object} ErrorResponse	"ID is incorrect"
// @Router      /api/posts/{id} [get]
func (h *Handler) GetPostById(c echo.Context) error {
	id, errReq := GetParam(c, ParamId)
	if errReq != nil {
		return errReq
	}

	post, err := h.services.Post.GetById(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "ID is incorrect.")
		return nil
	}
	_, errEnCd := json.Marshal(post)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, post)
	if errRes != nil {
		return errRes
	}
	return nil
}

// PostPost godoc
// @Summary      Add a post
// @Description  add post by json
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        post	body     PostRequest 		  true  "Add post"
// @Success      200 	{object} IdResponse		 "result is id of post"
// @Failure 	 400 	{object} ErrorResponse	 "user id is of valid type"
// @Failure 	 404 	{object} ErrorResponse	 "user id not found"
// @Failure 	 500 	{object} ErrorResponse	 "server error"
// @Router       /api/posts [post]
func (h *Handler) PostPost(c echo.Context) error {
	userId, errParams := GetUserId(c)
	if errParams != nil {
		return errParams
	}

	var post models.Post
	errReq := GetRequest(c, &post)
	if errReq != nil {
		return nil
	}

	post.UserId = userId
	id, err := h.services.Post.Create(post)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// UpdatePost godoc
// @Summary      Update a post
// @Description  Update by json post
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param       id  path     int true "Post ID"
// @Param        post	  body      PostRequest  true  "Update post"
// @Success      200      {object}  MessageResponse	"Post with id # updated"
// @Failure 	400 {object} ErrorResponse	 "incorrect request data"
// @Failure 	400 {object} ErrorResponse	 "user id is of valid type"
// @Failure 	400 {object} ErrorResponse	 "id is not integer"
// @Failure 	404 {object} ErrorResponse	 "user id not found"
// @Failure 	500 {object} ErrorResponse	 "server error"
// @Router       /api/posts/{id} [put]
func (h *Handler) UpdatePost(c echo.Context) error {
	id, errParams := GetParam(c, ParamId)
	if errParams != nil {
		return errParams
	}

	var post models.Post
	errReq := GetRequest(c, &post)
	if errReq != nil {
		return nil
	}

	err := h.services.Post.Update(id, post)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Post with id %d updated", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// DeletePost godoc
// @Summary      Delete a post
// @Description  Delete by json post
// @Tags         posts
// @Produce      json
// @Param       id  path     int true "Post ID"
// @Success     200 {object}  MessageResponse	"Post with id # deleted"
// @Failure 	400 {object} ErrorResponse	 "user id is of valid type"
// @Failure 	400 {object} ErrorResponse	 "id is not integer"
// @Failure 	404 {object} ErrorResponse	 "user id not found"
// @Failure 	500 {object} ErrorResponse	 "server error"
// @Router       /api/posts/{id} [delete]
func (h *Handler) DeletePost(c echo.Context) error {
	id, errParams := GetParam(c, ParamId)
	if errParams != nil {
		return nil
	}

	err := h.services.Post.Delete(id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	errRes := c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Post with id %d deleted", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
