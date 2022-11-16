package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"test"
)

func (h *Handler) SignUp(c echo.Context) error {
	var input test.User

	if errReq := c.Bind(&input); errReq != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	{
		//username is not empty
		if len(input.Username) == 0 {
			NewErrorResponse(c, http.StatusBadRequest, "You must enter a username")
			return nil
		}

		//name is not empty
		if len(input.Name) == 0 {
			NewErrorResponse(c, http.StatusBadRequest, "You must enter a name")
			return nil
		}

		// password length
		if len(input.Password) < 6 {
			NewErrorResponse(c, http.StatusBadRequest, "Password must be at least 6 symbols")
			return nil
		}
	}
	id, errCreate := h.services.Authorization.CreateUser(input)
	if errCreate != nil {
		NewErrorResponse(c, http.StatusInternalServerError, "something went wrong")
		return nil
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"id": id})
	if errRes != nil {
		return errRes
	}
	return nil
}

type SignInInput struct {
	Username string `json:"username" form:"username"  binding:"required"`
	Password string `json:"password" form:"password"  binding:"required"`
}

func (h *Handler) SignIn(c echo.Context) error {
	var input SignInInput
	if err := c.Bind(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect request data")
		return nil
	}
	errCheck := h.services.Authorization.CheckUser(input.Username)
	if errCheck != nil {
		NewErrorResponse(c, http.StatusBadRequest, "user not found")
		return nil
	}
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "incorrect password")
		return nil
	}
	errRes := c.JSON(http.StatusOK, map[string]interface{}{
		"token": token})
	if errRes != nil {
		return errRes
	}
	return nil
}

//func (h *Handler) SignInWithGoogle(c echo.Context) error {
//
//	// get request from google api
//	var input = test.User{
//		//Id:       0,
//		Name:     GetParam(c, Name),
//		Username: GetParam(c, Username),
//		Password: GetParam(c, Password),
//	}
//	fmt.Println(input)
//
//	//check username: if it is not used (it will get an error), try to create new user
//	errUser := h.services.Authorization.CheckUser(input.Username)
//
//	if errUser != nil {
//		h.GoogleSignUp(c, input)
//	}
//
//	// if username is required, to login and to generate token
//	h.GoogleSignIn(c, input)
//return nil
//}
//
//func (h Handler) Testing(c echo.Context) error {
//	var name string
//	errReq := GetRequest(c, name)
//	if errReq != nil {
//		return errReq
//	}
//
//	res, _ := h.services.Authorization.Testing(name)
//	c.JSON(http.StatusOK, res)
//	return nil
//}
