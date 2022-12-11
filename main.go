package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID 				uint			`json:"id,omitempty"`		
	Email 		string		`json:"email" binding:"required"`
	Password 	string		`json:"password,omitempty" binding:"required"`
}

var db = ConnectDB()

func main() {
	router := gin.Default()

	router.GET("/", Home)
	router.GET("/users", GetUsers)
	router.POST("/users", CreateUser)
	router.Run()
}

func Home(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"message": "Hello from Go.",
	})	
}

func GetUsers(c *gin.Context) {
	var user User
	var users []User
	rows, errQ := db.Query("SELECT * FROM users")
	if errQ != nil {
		panic(errQ)
	}
	for rows.Next() {
		errScan := rows.Scan(&user.ID, &user.Email, &user.Password)
		if errScan != nil {
			panic(errScan)
		}
		user = User{
			ID: user.ID,
			Email: user.Email,
		}
		users = append(users, user)
	}
	c.IndentedJSON(http.StatusOK, &users)
}

func CreateUser(c *gin.Context) {
	var newUser User

	if errBind := c.BindJSON(&newUser); errBind != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errBind.Error(),
		})
	} 
	result, errExec := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", newUser.Email, newUser.Password)
	if errExec != nil {
		return
	}
	insertID, _ := result.LastInsertId()
	c.IndentedJSON(http.StatusCreated, gin.H{
		"rows_affected": insertID,
	})
}