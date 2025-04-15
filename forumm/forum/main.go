package main

import (
    "github.com/gin-gonic/gin"
    "github.com/yourusername/forum/forum"
)

func main() {
    r := gin.Default()
    forumService := forum.NewService()
    r.GET("/posts", forumService.GetPosts)
    r.POST("/posts", forumService.CreatePost)
    r.Run(":8081")
}
