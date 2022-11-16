package handler

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"test"
	"test/pkg/service"
	mockService "test/pkg/service/mocks"
	"testing"
)

func TestHandler_userIdentify(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1" + "\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			e := echo.New()
			//e.Use(handler.userIdentify)
			e.GET("/protected", nil, handler.userIdentify, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					id := c.Get(userCtx).(int)
					c.JSON(200, id)
					return nil
				}
			})
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			e.ServeHTTP(rec, req)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}

func TestHandler_GoogleSignUp(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, user test.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            test.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name": "Test","username":"test username","password":"password"}`,
			inputUser: test.User{
				Name:     "Test",
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user test.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}` + "\n",
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "Test","username":"test username","password":"password"}`,
			inputUser: test.User{},
			mockBehavior: func(s *mockService.MockAuthorization, user test.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("something went wrong"))
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

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()
			//e.POST("/sign-up", handler.SignUp)

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPost, "/sign-up",
				//bytes.NewBufferString(testCase.inputBody))
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			//Проверка результатов
			if assert.NoError(t, handler.GoogleSignUp(ctx, testCase.inputUser)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func TestHandler_GoogleSignIn(t *testing.T) {
	type mockBehavior func(s *mockService.MockAuthorization, user test.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            test.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"username":"test username","password":"password"}`,
			inputUser: test.User{
				Username: "test username",
				Password: "password",
			},
			mockBehavior: func(s *mockService.MockAuthorization, user test.User) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("token", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"token"}` + "\n",
		},
		{
			name:      "Service Error",
			inputBody: `{"name": "Test","username":"test username","password":"password"}`,
			inputUser: test.User{},
			mockBehavior: func(s *mockService.MockAuthorization, user test.User) {
				s.EXPECT().GenerateToken(user.Username, user.Password).Return("", errors.New("something went wrong"))
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

			auth := mockService.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Тестовый сервер
			e := echo.New()
			//e.POST("/sign-up", handler.SignUp)

			//Тестовый запрос
			req := httptest.NewRequest(http.MethodPost, "/sign-in",
				//bytes.NewBufferString(testCase.inputBody))
				strings.NewReader(testCase.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			//Проверка результатов
			if assert.NoError(t, handler.GoogleSignIn(ctx, testCase.inputUser)) {
				assert.Equal(t, testCase.expectedStatusCode, rec.Code)
				assert.Equal(t, testCase.expectedResponseBody, rec.Body.String())
			}
		})
	}

}

func Test_GetParam(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(nil, rec)
	ctx.SetPath("/example/:value")
	ctx.SetParamNames("value")

	want := 4
	ctx.SetParamValues(fmt.Sprintf("%d", want))

	ok := struct {
		param int
		err   error
	}{}

	ok.param, ok.err = GetParam(ctx, "value")
	if ok.err != nil {
		t.Error("FAILED. Value of param dont have only numbers or null")
	} else if ok.param != want {
		t.Errorf("FAILED. Exepted %d, got %d", want, ok.param)
	} else {
		t.Logf("PASSED. Exepted %d, got %d", want, ok.param)
	}
}

func Test_GetRequest(t *testing.T) {

	inputBody := `{"name": "Test","username":"test username","password":"password"}`
	wantedUser := test.User{
		Name:     "Test",
		Username: "test username",
		Password: "password",
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/example",
		strings.NewReader(inputBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	var outputUser test.User

	err := GetRequest(ctx, &outputUser)
	if err != nil {
		t.Error("FAILED. Wrong request body")
	} else if assert.Equal(t, wantedUser, outputUser) {
		t.Logf("PASSED. Exepted %v, got %v", wantedUser, outputUser)
	} else {
		t.Errorf("FAILED. Exepted %v, got %v", wantedUser, outputUser)
	}
}

func Test_GetUserId(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	ctx := e.NewContext(nil, rec)
	want := 4
	ctx.Set(userCtx, want)

	ok := struct {
		param int
		err   error
	}{}

	ok.param, ok.err = GetUserId(ctx)
	if ok.err != nil {
		t.Error("FAILED. User not found or value is not a numbers")
	} else if ok.param != want {
		t.Errorf("FAILED. Exepted %d, got %d", want, ok.param)
	} else {
		t.Logf("PASSED. Exepted %d, got %d", want, ok.param)
	}
}
