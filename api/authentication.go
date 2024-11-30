package api

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
	"github.com/luideoz/API_Go/repository"
)

func Register(user models.User) int {
	status := repository.Add(user)
	return status
}

func Login(u models.User) int {
	user, status := repository.Get(u.USERNAME)
	if status != constants.OK {
		return status
	}
	if user.PASSWORD!= u.PASSWORD {
		return constants.INVALID_PASSWORD
	}
	return constants.OK
}

func Encrypt_hash(field string) string {
	hash := sha256.New()
	hash.Write([]byte(field))
	field = base64.StdEncoding.EncodeToString([]byte(field))
	return field
}
