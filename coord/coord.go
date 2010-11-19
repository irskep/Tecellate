package main

import (
	"fmt"
	"json"
	"os"
	"net"
	"log"
	"time"
)

type CoordConfig struct {
	identifier int
	BotConfs []BotConf
}

type BotConf struct {
	Path string
}

func main() {
	fmt.Printf("Launching with address: %s\n", os.Args[1])
	addr, err := net.ResolveTCPAddr(os.Args[1]);
	if err != nil { log.Exit(err) }
	listener, err := net.ListenTCP("tcp", addr);
	if err != nil { log.Exit(err) }
	
	conn, err := listener.AcceptTCP();
	if err != nil { log.Exit("error in Accept():", err) }
	defer conn.Close()
	conn.SetKeepAlive(true)
	conn.SetReadTimeout(30000)
	
	time.Sleep(1000000)
	
	fmt.Printf("Waiting for connections on %s...\n", addr)
	
	rcvd := make([]byte, 4096)
	size, err := conn.Read(rcvd)
	if err != nil { log.Exit("Error in recv:", err) }
	
	config := CoordConfig{0, nil}
	err = json.Unmarshal(rcvd, config)
	
	fmt.Printf("Received: %s\n", rcvd[0:size])
	
	conn.Write([]uint8("ok"))
}
