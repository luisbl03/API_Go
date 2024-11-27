package config

import (
	"log"
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
	err = os.Mkdir(Configs.Storage_root, 0755)
	if err != nil {
		log.Fatalf("Error creating root directory: %v", err)
	}
}
