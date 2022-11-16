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

func TestHandler_GetComments(t *testing.T) {
	type mockBehavior func(s *mockService.MockComment, postId int)

	testTable := []struct {
		name                 string
		paramId              int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "ok",
			paramId: 51,
			mockBehavior: func(s *mockService.MockComment, postId int) {
				ret := []test.Comment{
					{
						Id:     1,
						UserId: 20,
						Body:   "anons1",
					},
					{
						Id:     2,
						UserId: 31,
						Body:   "anons2",
					},
				}
				s.EXPECT().Get(postId).Return(ret, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"comments":[{"id":1,"user_id":20,"body":"anons1"},{"id":2,"user_id":31,"body":"anons2"}]}` + "\n",
		},
		{
			name:    "Server error",
			paramId: 51,
			mockBehavior: func(s *mockService.MockComment, postId int) {
				s.EXPECT().Get(postId).Return([]test.Comment{}, errors.New("something went wrong"))
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

			comment := mockService.NewMockComment(c)
			testCase.mockBehavior(comment, testCase.paramId)

			services := &service.Service{Comment: comment}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodGet, "/api/posts/:postId/comments", nil)
			//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/posts/:postId/comments")
			ctx.SetParamNames("postId")
			ctx.SetParamValues("51")

			//Проверка результатов
			if assert.NoError(t, handler.GetComments(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_PostComment(t *testing.T) {
	type mockBehavior func(s *mockService.MockComment, postId int, comment test.Comment)

	testTable := []struct {
		name                 string
		postId               int
		inputBody            string
		inputComment         test.Comment
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			postId:    3,
			inputBody: `{"body":"test body"}`,
			inputComment: test.Comment{
				UserId: 3,
				Body:   "test body",
			},
			mockBehavior: func(s *mockService.MockComment, postId int, comment test.Comment) {
				s.EXPECT().Create(postId, comment).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:      "server error",
			postId:    3,
			inputBody: `{"body":"test body"}`,
			inputComment: test.Comment{
				UserId: 3,
				Body:   "test body",
			},
			mockBehavior: func(s *mockService.MockComment, postId int, comment test.Comment) {
				s.EXPECT().Create(postId, comment).Return(0, errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		//{
		//	name:      "Error request data",
		//	inputBody: "error",
		//	inputPost: test.Post{
		//		Title: "test title",
		//		Anons: "test anons",
		//	},
		//	mockBehavior: func(s *mockService.MockComment, userId int, post test.Post) {
		//	},
		//	expectedStatusCode:   400,
		//	expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		//},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mockService.NewMockComment(c)
			testCase.mockBehavior(comment, testCase.postId, testCase.inputComment)

			services := &service.Service{Comment: comment}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPost, "/api/posts/:postId/comments",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(userCtx, 3)
			ctx.SetPath("/api/posts/:postId/comments")
			ctx.SetParamNames("postId")
			ctx.SetParamValues("3")

			//Проверка результатов
			if assert.NoError(t, handler.PostComment(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_UpdateComment(t *testing.T) {
	type mockBehavior func(s *mockService.MockComment, postId int, id int, comment test.Comment)

	testTable := []struct {
		name                 string
		postId               int
		commentId            int
		inputBody            string
		inputComment         test.Comment
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			postId:    3,
			commentId: 4,
			inputBody: `{"body":"test body"}`,
			inputComment: test.Comment{
				Body: "test body",
			},
			mockBehavior: func(s *mockService.MockComment, postId int, id int, comment test.Comment) {
				s.EXPECT().Update(postId, id, comment).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: `{"message":"Comment with id 4 updated."}` + "\n",
		},
		{
			name:      "server error",
			postId:    3,
			commentId: 4,
			inputBody: `{"body":"test body"}`,
			inputComment: test.Comment{
				Body: "test body",
			},
			mockBehavior: func(s *mockService.MockComment, postId int, id int, comment test.Comment) {
				s.EXPECT().Update(postId, id, comment).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		//{
		//	name:      "Error request data",
		//	inputBody: "error",
		//	userId:    3,
		//	postId:    1,
		//	inputPost: test.Post{
		//		Title: "test title",
		//		Anons: "test anons",
		//	},
		//	mockBehavior: func(s *mockService.MockComment, userId int, postId int, post test.Post) {
		//	},
		//	expectedStatusCode:   400,
		//	expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		//},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mockService.NewMockComment(c)
			testCase.mockBehavior(comment, testCase.postId, testCase.commentId, testCase.inputComment)

			services := &service.Service{Comment: comment}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPut, "/api/posts/:postId/comments/:id",
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Set(userCtx, 3)
			ctx.SetPath("/api/posts/:postId/comments/:id")
			ctx.SetParamNames("postId", "id")
			ctx.SetParamValues("3", "4")

			//Проверка результатов
			if assert.NoError(t, handler.UpdateComment(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_DeleteComment(t *testing.T) {
	type mockBehavior func(s *mockService.MockComment, postId int, id int)

	testTable := []struct {
		name                 string
		postId               int
		commentId            int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			postId:    3,
			commentId: 4,
			mockBehavior: func(s *mockService.MockComment, postId int, id int) {
				s.EXPECT().Delete(postId, id).Return(nil)
			},
			expectedStatusCode:   202,
			expectedResponseBody: `{"message":"Comment with id 4 deleted."}` + "\n",
		},
		{
			name:      "server error",
			postId:    3,
			commentId: 4,

			mockBehavior: func(s *mockService.MockComment, postId int, id int) {
				s.EXPECT().Delete(postId, id).Return(errors.New("server error"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"server error"}` + "\n",
		},
		//{
		//	name:      "Error request data",
		//	inputBody: "error",
		//	userId:    3,
		//	postId:    1,
		//	inputPost: test.Post{
		//		Title: "test title",
		//		Anons: "test anons",
		//	},
		//	mockBehavior: func(s *mockService.MockComment, userId int, postId int, post test.Post) {
		//	},
		//	expectedStatusCode:   400,
		//	expectedResponseBody: `{"message":"incorrect request data"}` + "\n",
		//},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {

			//Начальные значения
			//настраиваем логику оболочек (подключаем все уровни)
			c := gomock.NewController(t)
			defer c.Finish()

			comment := mockService.NewMockComment(c)
			testCase.mockBehavior(comment, testCase.postId, testCase.commentId)

			services := &service.Service{Comment: comment}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodDelete, "/api/posts/:postId/comments/:id", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.SetPath("/api/posts/:postId/comments/:id")
			ctx.SetParamNames("postId", "id")
			ctx.SetParamValues("3", "4")

			//Проверка результатов
			if assert.NoError(t, handler.DeleteComment(ctx)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}
