package models

import (
	"hello-from-go/config"
	"log"

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
		log.Panicln(err)
	}

	for rows.Next() {
		errScan := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
		if errScan != nil {
			log.Panicln(errScan)
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
		log.Panicln(errHash)
	}

	_, errExec := db.Exec("INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))", u.Email, string(hashP))
	return errExec == nil
}

func EmailExists(user User) bool {
	var email string
	checkEmail := db.QueryRow("SELECT email FROM users WHERE email = ?", user.Email)
	checkEmail.Scan(&email)
	return email != ""
}
