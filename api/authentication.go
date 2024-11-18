package api

import (
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/repository"
)

func Register(username string, password string) int {
	var user models.User
	models.SetUsername(&user, username)
	models.SetPassword(&user, password)
	status := repository.Add(user)
	return status
}

func Login(username string) (models.User, int) {
	user, status := repository.Get(username)
	return user, status
}
