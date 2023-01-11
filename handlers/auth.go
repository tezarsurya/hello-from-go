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
	"golang.org/x/crypto/bcrypt"
)

type JwtClaims struct {
	jwt.RegisteredClaims
	UserId  uint   `json:"userId"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var user models.User
	var cred Credentials

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

	password, _ := user.GetPassword(cred.Email)

	errCompare := bcrypt.CompareHashAndPassword([]byte(password), []byte(cred.Password))
	if errCompare != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid email or password"})
		return
	}

	errUser := user.FindByEmail(cred.Email)
	if errUser != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	claims := JwtClaims{
		jwt.RegisteredClaims{
			Issuer:    os.Getenv("BASE_URL"),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
		user.ID,
		user.Email,
		user.IsAdmin(),
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
