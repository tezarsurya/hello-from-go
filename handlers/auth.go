package handlers

import (
	"fmt"
	"hello-from-go/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId uint   `json:"user_id"`
	Email  string `json:"email"`
}

var cred models.Credentials

func Login(c *gin.Context) {
	if errBind := c.ShouldBindJSON(&cred); errBind != nil {
		var errValidation []gin.H
		for _, err := range errBind.(validator.ValidationErrors) {
			message := fmt.Sprintf("%s is required", err.Field())
			errValidation = append(errValidation, gin.H{
				"field": err.Field(),
				"error": message,
			})
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errValidation)
		return
	}
	id := cred.Login()
	if id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	claims := JwtClaims{
		jwt.RegisteredClaims{
			Issuer:    os.Getenv("BASE_URL"),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
		id,
		cred.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, errToken := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errToken != nil {
		log.Panicln(errToken)
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "logged in",
		"token":   signed,
	})
}
