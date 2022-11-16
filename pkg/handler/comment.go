package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"test"
)

type GetCommentsResponse struct {
	Comments []test.Comment `json:"comments"`
}

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
