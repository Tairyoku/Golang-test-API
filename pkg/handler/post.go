package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"test"
)

type GetPostsResponse struct {
	Posts []test.Post `json:"posts"`
}

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

func (h *Handler) GetUserPosts(c echo.Context) error {
	userId, errParams := GetParam(c, ParamId)
	if errParams != nil {
		NewErrorResponse(c, http.StatusBadRequest, "wrong param")
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

func (h *Handler) PostPost(c echo.Context) error {
	userId, errParams := GetUserId(c)
	if errParams != nil {
		return errParams
	}

	var post test.Post
	errReq := GetRequest(c, &post)
	if errReq != nil {
		return nil
	}

	id, err := h.services.Post.Create(userId, post)
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

func (h *Handler) UpdatePost(c echo.Context) error {
	id, errParams := GetParam(c, ParamId)
	if errParams != nil {
		return errParams
	}

	userId, errUserParams := GetUserId(c)
	if errUserParams != nil {
		return errUserParams
	}

	var post test.Post
	errReq := GetRequest(c, &post)
	if errReq != nil {
		return nil
	}

	err := h.services.Post.Update(userId, id, post)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}

	errRes := c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Post with id %d updated.", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}

func (h *Handler) DeletePost(c echo.Context) error {
	userId, errUserParams := GetUserId(c)
	if errUserParams != nil {
		return nil
	}

	id, errParams := GetParam(c, ParamId)
	if errParams != nil {
		return nil
	}

	err := h.services.Post.Delete(userId, id)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "server error")
		return nil
	}
	errRes := c.JSON(http.StatusAccepted, map[string]interface{}{
		"message": fmt.Sprintf("Post with id %d deleted.", id),
	})
	if errRes != nil {
		return errRes
	}
	return nil
}
