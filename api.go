package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int64  `json:"id"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
}

func main() {
	router := setupRouter()
	err := router.Run()

	fmt.Println(err)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.Use(ConnectDB())

	router.GET("/ping/:message", pingMessageController)
	router.POST("/initDB", initDBController)
	router.POST("/users", addUserController)
	router.GET("/users", listUserController)
	router.DELETE("/users/:id", deleteUserController)

	return router
}

func handleError(err error, code int, c *gin.Context) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": err.Error(),
	})
}

func ConnectDB() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := os.LookupEnv("DB_STRING")
		if !ok {
			val = "app.db"
		}

		database, err := sql.Open("sqlite3", val)
		if err != nil {
			handleError(err, 500, c)
			return
		}

		c.Set("db", database)
		c.Next()
		database.Close()
	}

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
		handleError(err, 500, c)
		return
	}
	_, err = statement.Exec()
	if err != nil {
		handleError(err, 500, c)
		return
	}

}

func addUserController(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	data, err := c.GetRawData()
	if err != nil {
		fmt.Println("a")
		handleError(err, 400, c)
		return
	}
	first_name, _, _, err := jsonparser.Get(data, "firstName")
	if err != nil {
		handleError(err, 400, c)
		return
	}
	last_name, _, _, err := jsonparser.Get(data, "lastName")
	if err != nil {
		handleError(err, 400, c)
		return
	}

	_, err = db.Exec("INSERT INTO people (firstname, lastname) VALUES (?, ?)", first_name, last_name)
	if err != nil {
		handleError(err, 500, c)
		return
	}
}

func listUserController(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	data, _ := db.Query("SELECT id, firstname, lastname from people")
	var users []User
	for data.Next() {
		var id int64
		var firstname string
		var lastname string
		_ = data.Scan(&id, &firstname, &lastname)
		users = append(users, User{id, firstname, lastname})
	}
	c.JSON(200, users)
}

func deleteUserController(c *gin.Context) {
	userId := c.Params.ByName("id")
	db := c.MustGet("db").(*sql.DB)

	_, err := db.Exec("DELETE FROM people where id = ?", userId)
	if err != nil {
		handleError(err, 500, c)
		return
	}

	c.JSON(204, nil)
}
