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
	router.POST("/initDB", initDBController)

	return router
}

func pingMessageController(c *gin.Context) {
	user := c.Params.ByName("message")

	c.JSON(200, gin.H{
		"message": "hello " + user,
	})
}

func initDBController(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	if err != nil {
		handleError(err, c)
		return
	}
	_, err = statement.Exec()
	if err != nil {
		handleError(err, c)
		return
	}

}

func ConnectDB() gin.HandlerFunc {

	return func(c *gin.Context) {
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
