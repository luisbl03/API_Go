package api

import (
	//"github.com/luideiz/API_Go/constants"
	"os"
	"path"
	"github.com/luideiz/API_Go/config"
	"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/repository"
)

func Upload(username string,json models.Json, id string) int {
	archive := id + ".json"
	path := path.Join(config.Configs.Storage_root, username, archive)
	status := repository.Upload(json, path)
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