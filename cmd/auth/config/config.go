package config

import (
	"os"

	"github.com/BurntSushi/toml"
)
type Config struct {
	Storage_root string `toml:"root"`
	Database_root string `toml:"database"`
}

var Configs Config

func Load(filename string){
	_, err := toml.DecodeFile(filename, &Configs)
	if err != nil {
		panic(err)
	}
	//miramos si existe el directorio
	_,err = os.Stat(Configs.Storage_root)
	if os.IsNotExist(err) {
		os.Mkdir(Configs.Storage_root, 0755)
	}
	_,err = os.Stat(Configs.Database_root)
	if os.IsNotExist(err) {
		os.Mkdir(Configs.Database_root, 0755)
		os.Create(Configs.Database_root+"/users.json")
	}
	_,err = os.Stat(Configs.Database_root+"/users.json")
	if os.IsNotExist(err) {
		os.Create(Configs.Database_root+"/users.json")
	}
}