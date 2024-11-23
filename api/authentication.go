package api

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/repository"
)

func Register(username string, password string) int {
	var user models.User
	hash := sha256.New()
	hash.Write([]byte(password))
	password = string(hash.Sum(nil))
	password_sha := base64.StdEncoding.EncodeToString([]byte(password))
	user.USERNAME = username
	user.PASSWORD = password_sha
	status := repository.Add(user)
	return status
}

func Login(username string, password string) (models.User, int) {
	user, status := repository.Get(username)
	if status != constants.OK {
		return user, status
	}
	password_user := user.PASSWORD
	hash := sha256.New()
	hash.Write([]byte(password_user))
	password_user = string(hash.Sum(nil))
	password_user = base64.StdEncoding.EncodeToString([]byte(password_user))
	if password_user != password {
		return user, constants.INVALID_PASSWORD
	}
	return user, constants.OK
}
