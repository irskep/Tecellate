package main

import (
	// "fmt"
	"os"
	"json"
	"io/ioutil"
	"log"
)

type Config struct {
	CoordPath string
}

var config = Config{""}

func main() {
	loadConfig()
	wd, err := os.Getwd()
	_, err = os.ForkExec(config.CoordPath, nil, nil, 
						 wd, []*os.File{nil, os.Stdout, os.Stderr})
	if err != nil {
		log.Exit(err)
	}
}

func loadConfig() {
	configFile, err := os.Open("config.json", os.O_RDONLY, 0)
	if err != nil {
		log.Exit(err)
	}
	defer configFile.Close()
	
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Exit(err)
	} else {
		json.Unmarshal(configBytes, &config)
	}
}
