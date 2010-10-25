package main

import (
	"fmt"
	"os"
	"json"
	"io/ioutil"
	"log"
	"strconv"
	"net"
	"time"
)

type Config struct {
	CoordPath string
	NumCoords int
	ListenAddr string
}

func main() {
	config := loadConfig()
	listenerReady := make(chan bool, 1)
	connections := make(map[int]*net.TCPConn)
	
	go listenForCoordinators(config, listenerReady, connections)
	
	fmt.Println("Launched getCoordinators, waiting for ok...")
	
	_ = <- listenerReady
	
	fmt.Println("Sweet, let's go")
	
	launchCoordinators(config)
	
	_ = <- listenerReady
}

func loadConfig() Config {
	var config = Config{"", 0, "127.0.0.1:4009"}
	
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
		argv := []string{strconv.Itoa(i), config.ListenAddr}
		_, err = os.ForkExec(config.CoordPath, argv, 
							 nil, wd, []*os.File{nil, os.Stdout, os.Stderr})
		if err != nil {
			log.Printf("Error on process %d:", i)
			log.Exit(err)
		}
	}
}

func listenForCoordinators(config Config, ready chan bool, connections map[int]*net.TCPConn) {
	addr, err := net.ResolveTCPAddr(config.ListenAddr);
	if err != nil { log.Exit(err) }
	listener, err := net.ListenTCP("tcp", addr);
	if err != nil { log.Exit(err) }
	
	fmt.Println("OK TO ISSUE")
	ready <- true
	
	for ; len(connections) < config.NumCoords; {
		conn, err := listener.AcceptTCP();
		if err != nil {
			log.Println("error in Accept():", err)
		} else {
			conn.SetKeepAlive(true)
			conn.SetReadTimeout(30000)
			rcvd := make([]byte, 127)
			size, err := conn.Read(rcvd)
			if err != nil { 
				if err == os.EAGAIN {
					log.Println("Recoverable?")
					for size, err = conn.Read(rcvd); err == os.EAGAIN;  {
						
					}
				} else {
					log.Println("Non-recoverable error")
					log.Exit(err)
				}
			}
			if size > 0 {
				numString := string(rcvd[0:size])
				id, err := strconv.Atoi(numString)
				if err != nil { log.Exit(err) }
				connections[id] = conn
			}
			time.Sleep(10000000)
		}
	}
	
	fmt.Println("Connected to all coordinators!")
	
	ready <- true
}
