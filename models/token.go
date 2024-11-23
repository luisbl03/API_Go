package models

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/luideiz/API_Go/constants"
)

// token struct
type Token struct {
	TOKEN string
	EXPIRATION time.Time
}

func IsAlive(token Token) bool {
	return token.EXPIRATION.After(time.Now())
}
func CreateToken() (Token, int) {
	var token Token
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