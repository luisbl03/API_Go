package repository

import (
	"encoding/json"
	"os"
	"github.com/luideiz/API_Go/constants"
	"github.com/luideiz/API_Go/models"
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