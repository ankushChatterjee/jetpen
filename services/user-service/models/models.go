package models

import (
	"encoding/hex"
	"log"
	"math/rand"
	"time"

	"github.com/ankushChatterjee/jetpen/user-service/pkg/utils"
)

type User struct {
	Username  string
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
}

type TempUser struct {
	Username  string
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
	Token     string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func NewTempUser(username string, email string, name string, password string) *TempUser {
	t := new(TempUser)
	t.Username = username
	t.Password = password
	t.Email = email
	t.Name = name
	t.Password, _ = utils.HashPassword(password)
	b := make([]byte, 56)
	if _, err := rand.Read(b); err != nil {
		failOnError(err, "Failed to create Token")
		return nil
	}
	t.Token = hex.EncodeToString(b)
	return t
}

func NewUser(username string, email string, name string, password string) *User {
	t := new(User)
	t.Username = username
	t.Password = password
	t.Email = email
	t.Name = name
	t.Password = password
	return t
}

func NewUserFromTempUser(tempUser *TempUser) *User {
	t := new(User)
	t.Username = tempUser.Username
	t.Password = tempUser.Password
	t.Email = tempUser.Email
	t.Name = tempUser.Name
	t.Password = tempUser.Password
	return t
}
