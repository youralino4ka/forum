package main

import (
    "github.com/gin-gonic/gin"
    "github.com/yourusername/forum/auth"
)

func main() {
    r := gin.Default()
    authService := auth.NewService()
    r.POST("/register", authService.Register)
    r.POST("/login", authService.Login)
    r.Run(":8080")
}