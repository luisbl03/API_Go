package api

import (
	//"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/repository"
	"crypto/md5"
	"encoding/hex"
)

func Upload(username string, data string, id string) int {
	var file models.Json
	hash := md5.New()
	hash.Write([]byte(id))
	file.Id = hex.EncodeToString(hash.Sum(nil))
	file.Username = username
	file.Data = data
	status:= repository.Add_json(file)
	return status
}