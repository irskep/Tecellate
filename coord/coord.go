package main

import (
	"fmt"
	"json"
	"os"
	"net"
	"log"
	"strconv"
	"time"
)

type CoordConfig struct {
	Identifier int
	BotConfs []BotConf
}

type BotConf struct {
	Path string
}

func dieIfError(err os.Error, msg string) {
	if err != nil { log.Exit("", msg, " in coordinator: ", err) }
}

func setupConnection() *net.TCPConn {
	fmt.Printf("Launching with address: %s\n", os.Args[1])
	addr, err := net.ResolveTCPAddr(os.Args[1]);
	dieIfError(err, "TCP address resolution error")
	listener, err := net.ListenTCP("tcp", addr);
	dieIfError(err, "Listening error")
	
	conn, err := listener.AcceptTCP();
	dieIfError(err, "TCP accept error")
	conn.SetKeepAlive(true)
	conn.SetReadTimeout(30000)
	
	fmt.Printf("Waiting for connections on %s...\n", addr)
	
	time.Sleep(1000000)
	
	return conn
}

func receive_from(conn *net.TCPConn) []byte {
	rcvd := make([]byte, 4096)
	size, err := conn.Read(rcvd)
	dieIfError(err, "Receive error")
	return rcvd[0:size]
}

func setupBot(conf BotConf, portNumber int) *net.TCPConn {
	addrString := "127.0.0.1:" + strconv.Itoa(portNumber)
	fd := []*os.File{os.Stdin,os.Stdout,os.Stderr};
    _, err := os.ForkExec(conf.Path, []string{addrString}, nil, "",fd);
	dieIfError(err, "Error launching bot")
	
	time.Sleep(100000000)
	
	addr, err := net.ResolveTCPAddr(addrString);
	dieIfError(err, "TCP address resolution error")
	conn, err := net.DialTCP("tcp", nil, addr)
	dieIfError(err, "Dial error")
	return conn
}

func setupBots(config *CoordConfig) []*net.TCPConn {
	connections := make([]*net.TCPConn, len(config.BotConfs))
	basePort := new(int)
	fmt.Sscanf(os.Args[1], "127.0.0.1:%d", basePort)
	// basePort, err := strconv.Atoi(portString)
	// fmt.Println("port: %d (%s)", basePort, portString)
	// dieIfError(err, "formatting fuck")
	for ix, b := range(config.BotConfs) {
		connections[ix] = setupBot(b, *basePort + ix + 1)
	}
	return connections
}

func main() {
	conn := setupConnection()
	defer conn.Close()
	
	config := new(CoordConfig)
	err := json.Unmarshal(receive_from(conn), config)
	dieIfError(err, "JSON error")
	
	connections := setupBots(config)
	
	for _, c := range(connections) {
		fmt.Printf("%s\n", receive_from(c))
	}
	
	conn.Write([]uint8("ok"))
}
