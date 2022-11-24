package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"test/pkg/repository/models"
	"test/pkg/service"
	mockService "test/pkg/service/mocks"
	"testing"
)

func TestHandler_GetPosts(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(s *mockService.MockPost) {
				ret := []models.Post{
					{
						Id:     1,
						UserId: 12,
						Title:  "title1",
						Anons:  "anons1",
					},
					{
						Id:     2,
						UserId: 15,
						Title:  "title2",
						Anons:  "anons2",
					},
				}
				s.EXPECT().Get().Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"posts":[{"id":1,"user_id":12,"title":"title1","anons":"anons1"},{"id":2,"user_id":15,"title":"title2","anons":"anons2"}]}` + "\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockService.MockPost) {
				s.EXPECT().Get().Return([]models.Post{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			post := mockService.NewMockPost(c)
			testCase.mockBehavior(post)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/api/posts", nil)
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			//Проверка результатов
			if assert.NoError(t, handler.GetPosts(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_GetUserPosts(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost, userId int)

	testTable := []struct {
		name                 string
		inputParam           int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "ok",
			inputParam: 12,
			mockBehavior: func(s *mockService.MockPost, userId int) {
				ret := []models.Post{
					{
						Id:     1,
						UserId: 12,
						Title:  "title1",
						Anons:  "anons1",
					},
					{
						Id:     2,
						UserId: 12,
						Title:  "title2",
						Anons:  "anons2",
					},
				}
				s.EXPECT().GetByUserId(userId).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"posts":[{"id":1,"user_id":12,"title":"title1","anons":"anons1"},{"id":2,"user_id":12,"title":"title2","anons":"anons2"}]}` + "\n",
		},
		{
			name:       "error param",
			inputParam: 12,
			mockBehavior: func(s *mockService.MockPost, userId int) {
				s.EXPECT().GetByUserId(userId).Return([]models.Post{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"wrong user ID"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			post := mockService.NewMockPost(c)
			testCase.mockBehavior(post, testCase.inputParam)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/api/posts/user/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/posts/user/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("12")

			//Проверка результатов
			if assert.NoError(t, handler.GetUserPosts(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_GetPostById(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost, id int)

	testTable := []struct {
		name                 string
		inputParam           int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:       "ok",
			inputParam: 1,
			mockBehavior: func(s *mockService.MockPost, id int) {
				ret := models.Post{
					Id:     1,
					UserId: 12,
					Title:  "title",
					Anons:  "anons",
				}
				s.EXPECT().GetById(id).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"user_id":12,"title":"title","anons":"anons"}` + "\n",
		},
		{
			name:       "error param",
			inputParam: 1,
			mockBehavior: func(s *mockService.MockPost, id int) {
				s.EXPECT().GetById(id).Return(models.Post{}, errors.New("ID is incorrect."))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"ID is incorrect."}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			post := mockService.NewMockPost(c)
			testCase.mockBehavior(post, testCase.inputParam)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/api/posts/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/posts/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			//Проверка результатов
			if assert.NoError(t, handler.GetPostById(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_PostPost(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost, post models.Post)

	testTable := []struct {
		name                 string
		paramUserId          int
		inputBody            string
		inputPost            models.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "ok",
			paramUserId: 12,
			inputBody:   `{"title":"test title","anons":"test anons"}`,
			inputPost: models.Post{
				UserId: 12,
				Title:  "test title",
				Anons:  "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, post models.Post) {
				s.EXPECT().Create(post).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:        "server error",
			paramUserId: 12,
			inputBody:   `{"title":"test title","anons":"test anons"}`,
			inputPost: models.Post{
				UserId: 12,
				Title:  "test title",
				Anons:  "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, post models.Post) {
				s.EXPECT().Create(post).Return(0, errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			inputPost: models.Post{
				UserId: 12,
				Title:  "test title",
				Anons:  "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, post models.Post) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			post := mockService.NewMockPost(c)
			testCase.mockBehavior(post, testCase.inputPost)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPost, "/api/posts",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(userCtx, 12)

			//Проверка результатов
			if assert.NoError(t, handler.PostPost(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_UpdatePost(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost, postId int, post models.Post)

	testTable := []struct {
		name                 string
		postId               int
		inputBody            string
		inputPost            models.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			postId:    1,
			inputBody: `{"title":"new title","anons":"new anons"}`,
			inputPost: models.Post{
				Title: "new title",
				Anons: "new anons",
			},
			mockBehavior: func(s *mockService.MockPost, postId int, post models.Post) {
				s.EXPECT().Update(postId, post).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: `{"message":"Post with id 1 updated"}` + "\n",
		},
		{
			name:      "server error",
			postId:    1,
			inputBody: `{"title":"test title","anons":"test anons"}`,
			inputPost: models.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, postId int, post models.Post) {
				s.EXPECT().Update(postId, post).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			postId:    1,
			inputPost: models.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, postId int, post models.Post) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			post := mockService.NewMockPost(c)
			testCase.mockBehavior(post, testCase.postId, testCase.inputPost)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPut, "/api/posts/:id",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/posts/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			//Проверка результатов
			if assert.NoError(t, handler.UpdatePost(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_DeletePost(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost, postId int)

	testTable := []struct {
		name                 string
		postId               int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "ok",
			postId: 1,
			mockBehavior: func(s *mockService.MockPost, postId int) {
				s.EXPECT().Delete(postId).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: `{"message":"Post with id 1 deleted"}` + "\n",
		},
		{
			name:   "server error",
			postId: 1,
			mockBehavior: func(s *mockService.MockPost, postId int) {
				s.EXPECT().Delete(postId).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			post := mockService.NewMockPost(c)
			testCase.mockBehavior(post, testCase.postId)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodDelete, "/api/posts/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/posts/:id")
			ctx.SetParamNames("id")
			ctx.SetParamValues("1")

			//Проверка результатов
			if assert.NoError(t, handler.DeletePost(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}
