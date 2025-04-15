package auth

import (
	"github.com/gin-gonic/gin"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Register(c *gin.Context) {
	// Логика регистрации пользователя
}

func (s *Service) Login(c *gin.Context) {
	// Логика авторизации пользователя
}
