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
 
func Add_json(file models.Json, username string) int {
	roots := []models.Folder{}
	jsonFile, err := os.Open(ARCHIVES)
	if err != nil {
		return constants.ERROR
	}
	defer jsonFile.Close()
	info,_ := jsonFile.Stat()
	before := info.Size()
	var after int64
	err = json.NewDecoder(jsonFile).Decode(&roots)
	if err != nil {
		return constants.ERROR
	}
	for i, root := range roots {
		if root.User == username {
			//miramos si existe el archivo en su directorio
			for _, f := range root.Files {
				if f.Id == file.Id {
					return constants.EXISTS
				}
			}
			roots[i].Files = append(roots[i].Files, file)
			jsonData, err := json.MarshalIndent(roots, "", "  ")
			if err != nil {
				return constants.ERROR
			}
			err = os.WriteFile(ARCHIVES, jsonData, 0644)
			if err != nil {
				return constants.ERROR
			}
			after = info.Size()
			break
		}
	}
	return int(after - before)
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