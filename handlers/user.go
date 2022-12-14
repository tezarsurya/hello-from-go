package handlers

import (
	"fmt"
	"hello-from-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetUsers(c *gin.Context) {
	var user models.User
	users := user.GetAll()
	if len(users) == 0 {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var user models.User
	if errBind := c.ShouldBindJSON(&user); errBind != nil {
		var errValidation []gin.H
		for _, err := range errBind.(validator.ValidationErrors) {
			var message string
			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("%s is required", err.Field())
			case "email":
				message = fmt.Sprintf("%s is invalid", err.Field())
			default:
				message = fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
			}
			errValidation = append(errValidation, gin.H{
				"field": err.Field(),
				"error": message,
			})
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, errValidation)
		return
	}

	if user.EmailExists() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "email already exists",
		})
		return
	}
	if !user.Create() {
		return
	}
	c.AbortWithStatus(http.StatusCreated)
}
