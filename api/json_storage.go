package api

import (
	//"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/repository"
)

func Upload(username string, data string, id string) int {
	var file models.Json
	file.Id = id
	file.Username = username
	file.Data = data
	status:= repository.Add_json(file)
	return status
}