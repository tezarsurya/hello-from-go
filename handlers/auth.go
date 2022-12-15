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
	if !cred.Login() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    cred.Email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
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
