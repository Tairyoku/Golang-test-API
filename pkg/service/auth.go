package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"test"
	"test/pkg/repository"
	"time"
)

const (
	salt      = "239tjeaWFYh2rofjw"
	tokenTTL  = 48 * time.Hour
	signInKey = "fl4i#kQgeg5leFrk&rkg43"
)

type AuthService struct {
	repository repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewAuthService(repository repository.Authorization) *AuthService {
	return &AuthService{repository: repository}
}

func (a *AuthService) Testing(name string) (string, error) {
	return a.repository.Testing(name)
}

func (a *AuthService) CreateUser(user test.User) (int, error) {
	user.Password = CreatePasswordHash(user.Password)
	return a.repository.CreateUser(user)
}

func (a *AuthService) CheckUser(username string) error {
	return a.repository.CheckUser(username)
}

func (a *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := a.repository.GetUser(username, CreatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	return token.SignedString([]byte(signInKey))
}

func (a *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signInKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserId, nil
}

func CreatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
