package models

import (
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (c *Credentials) Login() uint {
	user := User{
		ID:       0,
		Email:    c.Email,
		Password: c.Password,
	}
	password := db.QueryRow("SELECT id, password FROM users WHERE email = ?", c.Email)
	password.Scan(&user.ID, &user.Password)
	errCompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(c.Password))

	if !EmailExists(user) || errCompare != nil {
		return 0
	}
	return user.ID
}
