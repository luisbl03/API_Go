package repository

import (
	"encoding/json"
	"os"
	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)

const (
	ARCHIVES = "database/archives.json"
)
 
func Upload(data models.Json, path string) int {
	_,err := os.Stat(path)
	if err == nil {
		return constants.EXISTS
	}
	file,err := os.Create(path)
	if err != nil {
		return constants.ERROR
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return constants.ERROR
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return constants.ERROR
	}
	stat,err := os.Stat(path)
	if err != nil {
		return constants.ERROR
	}
	return int(stat.Size())
}

func GetFile(path string) (models.Json, int) {
	_,err := os.Stat(path)
	if err != nil {
		return models.Json{}, constants.NOT_FOUND
	}
	var data models.Json
	file,err := os.Open(path)
	if err != nil {
		return models.Json{}, constants.ERROR
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return models.Json{}, constants.ERROR
	}
	return data, constants.OK
}

func Update(path string, data models.Json) int {
	_,err := os.Stat(path)
	if err != nil {
		return constants.NOT_FOUND
	}
	//borrado del content y metemos lo nuevo
	file, err := os.Open(path)
	if err != nil {
		return constants.ERROR
	}
	defer file.Close()
	err = os.Truncate(path, 0)
	if err != nil {
		return constants.ERROR
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return constants.ERROR
	}
	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return constants.ERROR
	}
	stat,err := os.Stat(path)
	if err != nil {
		return constants.ERROR
	}
	size := stat.Size()
	return int(size)
}

func Delete (path string) int {
	_,err := os.Stat(path)
	if err != nil {
		return constants.NOT_FOUND
	}
	err = os.Remove(path)
	if err != nil {
		return constants.ERROR
	}
	return constants.OK
}

func List_Files(path string) ([]models.File, int) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, constants.ERROR
	}
	var files []models.File
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		var file models.Json
		file, status := GetFile(path + "/" + entry.Name())
		if status != constants.OK {
			return nil, status
		}
		files = append(files, models.File{Id: entry.Name(), Doc_content: file})
	}
	return files, constants.OK
}