package main

import (
	"fmt"
	"os"
	"json"
	"io/ioutil"
	"log"
	"strconv"
)

type Config struct {
	CoordPath string
}

var config = Config{""}

func main() {
	loadConfig()
	ch := make(chan int)
	go getCoordinators(ch)
	fmt.Println("Launched getCoordinators, waiting for ok...")
	_ = <-ch
	fmt.Println("Sweet, let's go")
	launchCoordinators()
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

func launchCoordinators() {
	wd, err := os.Getwd()
	for i := 0; i < 10; i++ {
		_, err = os.ForkExec(config.CoordPath, []string{strconv.Itoa(i)}, nil, 
							 wd, []*os.File{nil, os.Stdout, os.Stderr})
		if err != nil {
			log.Printf("Error on process %d:", i)
			log.Exit(err)
		}
	}
}

func getCoordinators(ok chan int) {
	fmt.Println("OK TO ISSUE")
	ok <- 1
}
