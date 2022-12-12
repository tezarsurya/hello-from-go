package models

import (
	"database/sql"
	"hello-from-go/config"
)

type User struct {
	ID 				uint			`json:"id,omitempty"`		
	Email 		string		`json:"email" binding:"required"`
	Password 	string		`json:"password,omitempty" binding:"required"`
}

var db = config.ConnectDB()

func GetAllUsers() []User {
	var user User
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		errScan := rows.Scan(&user.ID, &user.Email, &user.Password)
		if errScan != nil {
			panic(errScan)
		}
		user = User{
			Email: user.Email,
			Password: user.Password,
		}
		users = append(users, user)
	}

	return users
}

func NewUser(newUser User) (sql.Result, error) {
	result, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", newUser.Email, newUser.Password)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func CheckEmail(user User) bool {
	var email string
	checkEmail := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email)
	checkEmail.Scan(&email)
	return email == ""
}
