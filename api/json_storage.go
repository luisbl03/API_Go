package api

import (
	//"github.com/luideiz/API_Go/constants"
	"os"
	"path"
	"github.com/luideoz/API_Go/config"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
	"github.com/luideoz/API_Go/repository"
)

func Upload(username string,json models.Json, id string) int {
	archive := id + ".json"
	path := path.Join(config.Configs.Storage_root, username, archive)
	//si no existe el directorio, lo crea
	if _, err := os.Stat(config.Configs.Storage_root + "/" + username); os.IsNotExist(err) {
		status := Root(username)
		if status != constants.OK {
			return status
		}
	}
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

func GetFile(id string, user string) (models.Json, int) {
	archive := id + ".json"
	path := path.Join(config.Configs.Storage_root, user, archive)
	json, status := repository.GetFile(path)
	return json, status
}

func Update(id string, user string, json models.Json) int {
	archive := id + ".json"
	path := path.Join(config.Configs.Storage_root, user, archive)
	status := repository.Update(path, json)
	return status
}

func Delete(id string, username string) int {
	archive := id + ".json"
	path := path.Join(config.Configs.Storage_root, username, archive)
	status := repository.Delete(path)
	return status
}

func List_Files(user string) ([]models.File, int) {
	path := config.Configs.Storage_root + "/" + user
	files, status := repository.List_Files(path)
	return files, status
}