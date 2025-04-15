package forum

import (
	"github.com/gin-gonic/gin"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetPosts(c *gin.Context) {
	// Логика получения постов
}

func (s *Service) CreatePost(c *gin.Context) {
	// Логика создания поста
}
