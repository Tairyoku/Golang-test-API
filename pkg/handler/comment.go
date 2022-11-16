package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"test"
)

// GetComments godoc
// @Summary     Find all comments
// @Description Get all comments
// @Tags        comments
// @Produce     json
// @Param       postId  path     int true "Post ID"
// @Success     200 {object} GetCommentsResponse
// @Failure 	400 {object} ErrorResponse	 "postId is not integer"
// @Failure 	500 {object} ErrorResponse	 "something went wrong"
// @Router      /api/posts/{postId}/comments [get]
func (h *Handler) GetComments(c echo.Context) error {
	postId, errParams := GetParam(c, ParamPostId)
	if errParams != nil {
		return errParams
	}

	comments, err := h.services.Comment.Get(postId)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "something went wrong")
		return nil
	}
	_, errEnCd := json.Marshal(comments)
	if errEnCd != nil {
		return errEnCd
	}
	errRes := c.JSON(http.StatusOK, GetCommentsResponse{
		Comments: comments,
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// PostComment godoc
// @Summary      Add a comment
// @Description  add by json comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param        post	body     CommentRequest   true  "Add comment"
// @Success      200 	{object} IdResponse		 "result is id of comment"
// @Failure 	 400 	{object} ErrorResponse	 "postId not integer"
// @Failure 	 400 	{object} ErrorResponse	 "incorrect request data"
// @Failure 	 400 	{object} ErrorResponse	 "user id is of valid type"
// @Failure 	 404 	{object} ErrorResponse	 "user id not found"
// @Failure 	 500 	{object} ErrorResponse	 "server error"
// @Router       /api/posts/{postId}/comments [post]
func (h *Handler) PostComment(c echo.Context) error {
	postId, errParams := GetParam(c, ParamPostId)
	if errParams != nil {
		return errParams
	}

	userId, errUserParams := GetUserId(c)
	if errUserParams != nil {
		return nil
	}
	var comment test.Comment
	errReq := GetRequest(c, &comment)
	if errReq != nil {
		return errReq
	}

	comment.UserId = userId

	id, err := h.services.Comment.Create(postId, comment)
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

// UpdateComment godoc
// @Summary      Update a comment
// @Description  Update by json comment
// @Tags         comments
// @Accept       json
// @Produce      json
// @Param       id  path     int true "Comment ID"
// @Param       postId  path     int true "Post ID"
// @Param        post	  body      CommentRequest  true  "Update comment"
// @Success      200      {object}  MessageResponse	"Comment with id # updated"
// @Failure 	400 {object} ErrorResponse	 "incorrect request data"
// @Failure 	400 {object} ErrorResponse	 "id is not integer"
// @Failure 	400 {object} ErrorResponse	 "postId is not integer"
// @Failure 	500 {object} ErrorResponse	 "server error"
// @Router       /api/posts/{postId}/comments/{id} [put]
func (h *Handler) UpdateComment(c echo.Context) error {
	id, errParamId := GetParam(c, ParamId)
	if errParamId != nil {
		return errParamId
	}

	postId, errParams := GetParam(c, ParamPostId)
	if errParams != nil {
		return errParams
	}

	var comment test.Comment
	errReq := GetRequest(c, &comment)
	if errReq != nil {
		return errReq
	}

	err := h.services.Comment.Update(postId, id, comment)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Comment with id %d updated.", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Delete by json comment
// @Tags         comments
// @Produce      json
// @Param       id  path     int true "Post ID"
// @Param       postId  path     int true "Post ID"
// @Success     200 {object}  MessageResponse	"comment with id # deleted"
// @Failure 	400 {object} ErrorResponse	 "id is not integer"
// @Failure 	400 {object} ErrorResponse	 "postId is not integer"
// @Failure 	500 {object} ErrorResponse	 "server error"
// @Router       /api/posts/{postId}/comments/{id} [delete]
func (h *Handler) DeleteComment(c echo.Context) error {
	id, errParamId := GetParam(c, ParamId)
	if errParamId != nil {
		return errParamId
	}

	postId, errParams := GetParam(c, ParamPostId)
	if errParams != nil {
		return errParams
	}

	err := h.services.Comment.Delete(postId, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	errRes := c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Comment with id %d deleted.", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
