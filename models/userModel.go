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

func (u *User) GetAll() []User {
	db := config.ConnectDB()
	var users []User

	rows, err := db.Query("SELECT id, email, created_at, updated_at, deleted_at FROM users")
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		errScan := rows.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
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

func (u *User) FindByEmail(email string) error {
	db := config.ConnectDB()
	query := db.QueryRow("SELECT id, email, created_at, updated_at, deleted_at FROM users WHERE email = ?", email)
	err := query.Scan(&u.ID, &u.Email, &u.CreatedAt, &u.UpdatedAt, &u.DeletedAt)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetPassword(email string) (string, error) {
	db := config.ConnectDB()
	query := db.QueryRow("SELECT password from users WHERE email = ?", email)
	err := query.Scan(&u.Password)
	if err != nil {
		return "", err
	}
	return u.Password, nil
}

func (u *User) Create() bool {
	db := config.ConnectDB()
	hashP, errHash := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if errHash != nil {
		log.Panicln(errHash)
	}

	_, errExec := db.Exec("INSERT INTO users (email, password, created_at, updated_at) VALUES (?, ?, datetime('now'), datetime('now'))", u.Email, string(hashP))
	return errExec == nil
}

func (u *User) IsAdmin() bool {
	db := config.ConnectDB()
	var email string
	usrdata := db.QueryRow("SELECT email FROM admins WHERE email = ?", u.Email)
	err := usrdata.Scan(&email)

	if email == "" || err != nil {
		return false
	}
	return true
}

func (u *User) EmailExists() bool {
	db := config.ConnectDB()
	var email string
	checkEmail := db.QueryRow("SELECT email FROM users WHERE email = ?", u.Email)
	checkEmail.Scan(&email)
	return email != ""
}
