package config

import (
	"os"

	"github.com/BurntSushi/toml"
)
type Config struct {
	Storage_root string `toml:"root"`
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
}