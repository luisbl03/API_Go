package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	"github.com/luideoz/API_Go/constants"
)

// token struct
type Token struct {
	TOKEN string
	EXPIRATION time.Time
	User string
}

func IsAlive(token Token) bool {
	return token.EXPIRATION.After(time.Now())
}
func CreateToken(user string) (Token, int) {
	var token Token
	token.User = user
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return token, constants.ERROR
	}
	token.TOKEN = base64.StdEncoding.EncodeToString(bytes)
	// le ponemos de duracion 2 minutos
	token.EXPIRATION = time.Now().Add(time.Minute * 2)
	return token, constants.OK
}