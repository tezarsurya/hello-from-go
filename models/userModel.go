package models

import (
	"database/sql"
	"hello-from-go/config"
)

type User struct {
	ID       uint   `json:"id,omitempty"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password,omitempty" binding:"required,min=8"`
}

var db = config.ConnectDB()

func (u *User) GetAll() []User {
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		errScan := rows.Scan(&u.ID, &u.Email, &u.Password)
		if errScan != nil {
			panic(errScan)
		}
		*u = User{
			Email:    u.Email,
			Password: u.Password,
		}
		users = append(users, *u)
	}

	return users
}

func (u *User) Create() sql.Result {
	result, err := db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", u.Email, u.Password)
	if err != nil {
		panic(err)
	}
	return result
}

func CheckEmail(user User) bool {
	var email string
	checkEmail := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email)
	checkEmail.Scan(&email)
	return email == ""
}
