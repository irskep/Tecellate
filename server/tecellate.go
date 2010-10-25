package main

import (
    "fmt"
    "os"
    "json"
    "io/ioutil"
    "log"
)

type Config struct {
    CoordPath string
}

func main() {
    fmt.Printf("hello\n")
	var config = Config{""}
	
	configFile, err := os.Open("config.json", os.O_RDONLY, 0)
	if err != nil {
        return
    }
    defer configFile.Close()
    
    configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
        log.Exit(err)
    } else {
        json.Unmarshal(configBytes, &config)
    }
    fmt.Printf("%s\n", config.CoordPath)
    fmt.Printf("goodbye\n")
}
