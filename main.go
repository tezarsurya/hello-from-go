package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", Home)
	router.Run()
}

func Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Hello from Go.",
	})
}