package repository

import (
	"encoding/json"
	"os"

	"github.com/luideoz/API_Go/constants"
	"github.com/luideoz/API_Go/models"
)
const (
	FILE = "database/users.json"
) 


func Add(user models.User) int {
	// Open our jsonFile
	// 0 -> todo correcto
	// 1 -> error por usuario ya existente
	// 2 -> otro error
	usuarios := []models.User{}
	jsonFile, err := os.Open(FILE)
	if err != nil {
		return constants.ERROR
	}
	defer jsonFile.Close()
	//miramos si el archivo esta vacio
	info, _ := jsonFile.Stat()
	if info.Size() == 0 {
		usuarios = append(usuarios, user)
		jsonData, err := json.MarshalIndent(usuarios, "", "  ")
		if err != nil {
			return constants.ERROR
		}
		err = os.WriteFile(FILE, jsonData, 0644)
		if err != nil {
			return constants.ERROR
		}
		return constants.OK
	}
	//miramos si el usuario ya existe
	err = json.NewDecoder(jsonFile).Decode(&usuarios)
	if err != nil {
		return constants.ERROR
	}
	for _, u := range usuarios {
		if u.USERNAME == user.USERNAME {
			return constants.EXISTS
		}
	}
	usuarios = append(usuarios, user)
	jsonData, err := json.MarshalIndent(usuarios, "", "  ")
	if err != nil {
		return constants.ERROR
	}
	err = os.WriteFile(FILE, jsonData, 0644)
	if err != nil {
		return constants.ERROR
	}
	return constants.CREATED
}

func Get(username string) (models.User, int) {
	usuarios := []models.User{}
	jsonFile, err := os.Open(FILE)
	if err != nil {
		return models.User{}, constants.ERROR
	}
	defer jsonFile.Close()
	//miramos si el usuario ya existe
	err = json.NewDecoder(jsonFile).Decode(&usuarios)
	if err != nil {
		return models.User{}, constants.ERROR
	}
	for _, u := range usuarios {
		if u.USERNAME == username{
			return u, constants.OK
		}
	}
	return models.User{}, constants.NOT_FOUND
}