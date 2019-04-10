package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create gin application
	router := setupRouter()

	// Launch the application
	err := router.Run()

	fmt.Println(err)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(ConnectDB())

	router.GET("/ping/:message", pingMessageController)

	return router
}

func pingMessageController(c *gin.Context) {
	user := c.Params.ByName("message")

	c.JSON(200, gin.H{
		"message": "hello " + user,
	})
}

func ConnectDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		os.LookupEnv("DB_STRING")
		val, ok := os.LookupEnv("DB_STRING")
		if !ok {
			val = "app.db"
		}

		database, err := sql.Open("sqlite3", val)

		if err != nil {
			handleError(err, c)
			return
		}

		c.Set("db", database)

		c.Next()

		database.Close()
	}

}

func handleError(err error, c *gin.Context) {
	c.AbortWithStatusJSON(500, gin.H{
		"error": err.Error(),
	})
}
