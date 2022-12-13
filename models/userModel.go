package models

import (
	"hello-from-go/config"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint        `json:"id,omitempty"`
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password,omitempty" binding:"required,min=8"`
	CreatedAt string      `json:"createdAt"`
	UpdatedAt string      `json:"updatedAt,omitempty"`
	DeletedAt interface{} `json:"deletedAt,omitempty"`
}

var db = config.ConnectDB()

func (u *User) GetAll() []User {
	var users []User

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		errScan := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
		if errScan != nil {
			panic(errScan)
		}
		*u = User{
			ID:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
			DeletedAt: u.DeletedAt,
		}
		users = append(users, *u)
	}
	return users
}

func (u *User) Create() bool {
	hashP, errHash := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if errHash != nil {
		panic(errHash)
	}

	created := time.Now()
	_, errExec := db.Exec("INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, ?, ?)", u.Email, string(hashP), created, created)
	return errExec == nil
}

func CheckEmail(user User) bool {
	var email string
	checkEmail := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email)
	checkEmail.Scan(&email)
	return email == ""
}
