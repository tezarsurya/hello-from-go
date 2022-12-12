package main

import (
	"hello-from-go/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	router := gin.Default()

	router.GET("/", Home)
	router.GET("/users", handlers.GetUsers)
	router.POST("/users", handlers.CreateUser)
	router.Run()
}

func Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Hello from Go.",
	})	
}
