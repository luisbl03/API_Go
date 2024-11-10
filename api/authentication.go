package api

import (
	"github.com/luideiz/API_Go/models"
)

func Login(username string, password string) bool {
	var user models.User
	models.SetUsername(&user, username)
	models.SetPassword(&user, password)

	
	return true
}
