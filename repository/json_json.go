package repository

import (
	"encoding/json"
	"github.com/luideiz/API_Go/models"
	"github.com/luideiz/API_Go/constants"
	"os"
)

const (
	ARCHIVES = "database/archives.json"
)
 
func Add_json(file models.Json) int {
	archivos := []models.Json{}
	jsonFile, err := os.Open(ARCHIVES)
	if err != nil {
		return constants.ERROR
	}
	defer jsonFile.Close()
	/*miramos si el archivo esta vacio*/
	fileInfo,_ := os.Stat(ARCHIVES)
	before := fileInfo.Size()
	info, _ := jsonFile.Stat()
	if info.Size() == 0 {
		archivos = append(archivos, file)
	} else {
		/*miramos si no existe el json*/
		err = json.NewDecoder(jsonFile).Decode(&archivos)
		if err != nil {
			return constants.ERROR
		}
		for _, u := range archivos {
			if u.Id == file.Id && u.Username == file.Username {
				return constants.EXISTS
			}
		}
		archivos = append(archivos, file)
	}
	jsonData, err := json.MarshalIndent(archivos, "", "  ")
	if err != nil {
		return constants.ERROR
	}
	err = os.WriteFile(ARCHIVES, jsonData, 0644)
	if err != nil {
		return constants.ERROR
	}
	fileInfo,_ = os.Stat(ARCHIVES)
	after := fileInfo.Size()
	return int(after-before)
}

func Root(folder models.Folder) int {
	jsonFile, err := os.Open(ARCHIVES)
	if err != nil {
		return constants.ERROR
	}
	defer jsonFile.Close()
	roots := []models.Folder{}
	info, _ := jsonFile.Stat()
	if info.Size() == 0 {
		roots = append(roots, folder)
	} else {
		err = json.NewDecoder(jsonFile).Decode(&roots)
		if err != nil {
			return constants.ERROR
		}
		roots = append(roots, folder)
	}
	jsonData, err := json.MarshalIndent(roots, "", "  ")
	if err != nil {
		return constants.ERROR
	}
	err = os.WriteFile(ARCHIVES, jsonData, 0644)
	if err != nil {
		return constants.ERROR
	}
	return constants.OK
}