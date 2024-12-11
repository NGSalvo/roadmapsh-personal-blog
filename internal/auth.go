package internal

import (
	"golang.org/x/crypto/bcrypt"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

var users = map[string]Login{}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func AddUser(username string, password string) {
	users[username] = Login{
		HashedPassword: password,
	}
}

func GetUser(username string) *Login {
	user, ok := users[username]

	if !ok {
		return nil
	}

	return &user
}
