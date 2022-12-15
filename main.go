package main

import (
	"hello-from-go/handlers"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	godotenv.Load()
	router := gin.Default()

	router.GET("/", Home)
	router.POST("/login", handlers.Login)

	users := router.Group("users", AuthRequired())
	{
		users.GET("", handlers.GetUsers)
		users.POST("", handlers.CreateUser)
	}
	router.Run()
}

func Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Hello from Go.",
	})
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		reg := regexp.MustCompile(`Bearer [a-zA-z0-9]+[.][a-zA-z0-9]+[.][a-zA-z0-9-]+`)
		token := c.GetHeader("Authorization")
		if !reg.MatchString(token) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token = token[7:]
		parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods(jwt.GetAlgorithms()))

		if !parsed.Valid || err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
