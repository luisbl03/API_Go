package api

import (
	//"github.com/luideiz/API_Go/constants"
	"crypto/md5"
	"encoding/hex"
	"os"

	"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/repository"
	"github.com/luideiz/API_Go/config"
)

func Upload(username string, data string, id string) int {
	var file models.Json
	hash := md5.New()
	hash.Write([]byte(id))
	file.Id = hex.EncodeToString(hash.Sum(nil))
	file.Data = data
	status:= repository.Add_json(file, username)
	return status
}

func Root(username string) int {
	path := config.Configs.Storage_root + "/" + username
	err := os.Mkdir(path, 0755)
	if err != nil {
		return constants.ERROR
	}
	return constants.OK
}