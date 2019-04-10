package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := setupRouter()
	err := router.Run()

	fmt.Println(err)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/ping/:message", pingMessageController)

	return router
}

func pingMessageController(c *gin.Context) {
	user := c.Params.ByName("message")
	c.JSON(200, gin.H{
		"message": "hello " + user,
	})
}
