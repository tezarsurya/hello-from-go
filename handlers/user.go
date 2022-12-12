package handlers

import (
	"hello-from-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	users := models.GetAllUsers()
	c.IndentedJSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	
	if errBind := c.BindJSON(&newUser); errBind != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errBind.Error(),
		})
		return
	} 
			
	if !models.CheckEmail(newUser) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "email already exists",
		})
		return
	}

	result, errInsert := models.NewUser(newUser)
	if errInsert != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": errInsert.Error(),
		})
		return
	}
	insertID, _ := result.LastInsertId()
	c.IndentedJSON(http.StatusCreated, gin.H{
		"inserted id": insertID,
	})
}