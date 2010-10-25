package main

import (
	"fmt"
	"os"
	"json"
	"io/ioutil"
	"log"
	"strconv"
	"net"
)

type Config struct {
	CoordPath string
	NumCoords uint
	ListenPort uint
}

func main() {
	config := loadConfig()
	listenerReady := make(chan bool, 1)
	connections := map[int]*net.TCPConn{}
	
	go listenForCoordinators(config, listenerReady, connections)
	
	fmt.Println("Launched getCoordinators, waiting for ok...")
	
	_ = <- listenerReady
	
	fmt.Println("Sweet, let's go")
	
	launchCoordinators(config)
}

func loadConfig() Config {
	var config = Config{"", 0, 4009}
	
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
	return config
}

func launchCoordinators(config Config) {
	wd, err := os.Getwd()
	for i := 0; i < 10; i++ {
		argv := []string{strconv.Itoa(i), strconv.Uitoa(config.ListenPort)}
		_, err = os.ForkExec(config.CoordPath, argv, 
							 nil, wd, []*os.File{nil, os.Stdout, os.Stderr})
		if err != nil {
			log.Printf("Error on process %d:", i)
			log.Exit(err)
		}
	}
}

func listenForCoordinators(config Config, ready chan bool, connections map[int]*net.TCPConn) {
	addr, err := net.ResolveTCPAddr("127.0.0.1:4009");
	if err != nil { log.Exit(err) }
	listener, err := net.ListenTCP("tcp", addr);
	if err != nil { log.Exit(err) }
	
	for {
		conn, err := listener.AcceptTCP();
		if err != nil {
			log.Println("error in Accept():", err)
		} else {
			conn.SetKeepAlive(true);
			conn.SetReadTimeout(30000);
		}
	}
	
	fmt.Println("OK TO ISSUE")
	ready <- true
}
