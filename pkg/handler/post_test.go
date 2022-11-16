package handler

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"test"
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
				ret := []test.Post{
					{
						Id:    1,
						Title: "title1",
						Anons: "anons1",
					},
					{
						Id:    2,
						Title: "title2",
						Anons: "anons2",
					},
				}
				s.EXPECT().Get().Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"posts":[{"id":1,"title":"title1","anons":"anons1"},{"id":2,"title":"title2","anons":"anons2"}]}` + "\n",
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockService.MockPost) {
				s.EXPECT().Get().Return([]test.Post{}, errors.New("something went wrong"))
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
			inputParam: 2,
			mockBehavior: func(s *mockService.MockPost, userId int) {
				ret := []test.Post{
					{
						Id:    1,
						Title: "title1",
						Anons: "anons1",
					},
					{
						Id:    2,
						Title: "title2",
						Anons: "anons2",
					},
				}
				s.EXPECT().GetByUserId(userId).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"posts":[{"id":1,"title":"title1","anons":"anons1"},{"id":2,"title":"title2","anons":"anons2"}]}` + "\n",
		},
		{
			name:       "error param",
			inputParam: 2,
			mockBehavior: func(s *mockService.MockPost, userId int) {
				s.EXPECT().GetByUserId(userId).Return([]test.Post{}, errors.New("something went wrong"))
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
			ctx.SetParamValues("2")

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
				ret := test.Post{
					Id:    1,
					Title: "title",
					Anons: "anons",
				}
				s.EXPECT().GetById(id).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1,"title":"title","anons":"anons"}` + "\n",
		},
		{
			name:       "error param",
			inputParam: 1,
			mockBehavior: func(s *mockService.MockPost, id int) {
				s.EXPECT().GetById(id).Return(test.Post{}, errors.New("ID is incorrect."))
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
	type mockBehavior func(s *mockService.MockPost, userId int, post test.Post)

	testTable := []struct {
		name                 string
		userId               int
		inputBody            string
		inputPost            test.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			userId:    3,
			inputBody: `{"title":"test title","anons":"test anons"}`,
			inputPost: test.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, userId int, post test.Post) {
				s.EXPECT().Create(userId, post).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:      "server error",
			userId:    3,
			inputBody: `{"title":"test title","anons":"test anons"}`,
			inputPost: test.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, userId int, post test.Post) {
				s.EXPECT().Create(userId, post).Return(0, errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			inputPost: test.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, userId int, post test.Post) {
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
			testCase.mockBehavior(post, testCase.userId, testCase.inputPost)

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
			ctx.Set(userCtx, 3)

			//Проверка результатов
			if assert.NoError(t, handler.PostPost(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_UpdatePost(t *testing.T) {
	type mockBehavior func(s *mockService.MockPost, userId int, postId int, post test.Post)

	testTable := []struct {
		name                 string
		userId               int
		postId               int
		inputBody            string
		inputPost            test.Post
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			userId:    3,
			postId:    1,
			inputBody: `{"title":"new title","anons":"new anons"}`,
			inputPost: test.Post{
				Title: "new title",
				Anons: "new anons",
			},
			mockBehavior: func(s *mockService.MockPost, userId int, postId int, post test.Post) {
				s.EXPECT().Update(userId, postId, post).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: `{"message":"Post with id 1 updated."}` + "\n",
		},
		{
			name:      "server error",
			userId:    3,
			postId:    1,
			inputBody: `{"title":"test title","anons":"test anons"}`,
			inputPost: test.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, userId int, postId int, post test.Post) {
				s.EXPECT().Update(userId, postId, post).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		{
			name:      "Error request data",
			inputBody: "error",
			userId:    3,
			postId:    1,
			inputPost: test.Post{
				Title: "test title",
				Anons: "test anons",
			},
			mockBehavior: func(s *mockService.MockPost, userId int, postId int, post test.Post) {
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
			testCase.mockBehavior(post, testCase.userId, testCase.postId, testCase.inputPost)

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
			ctx.Set(userCtx, 3)
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
	type mockBehavior func(s *mockService.MockPost, userId int, postId int)

	testTable := []struct {
		name                 string
		userId               int
		postId               int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:   "ok",
			userId: 3,
			postId: 1,
			mockBehavior: func(s *mockService.MockPost, userId int, postId int) {
				s.EXPECT().Delete(userId, postId).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: `{"message":"Post with id 1 deleted."}` + "\n",
		},
		{
			name:   "server error",
			userId: 3,
			postId: 1,

			mockBehavior: func(s *mockService.MockPost, userId int, postId int) {
				s.EXPECT().Delete(userId, postId).Return(errors.New("server error"))
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
			testCase.mockBehavior(post, testCase.userId, testCase.postId)

			services := &service.Service{Post: post}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodDelete, "/api/posts/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(userCtx, 3)
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
