package main

import (
	"fmt"
	"json"
	"os"
	"net"
	"strconv"
	"strings"
	"time"
)

type CoordConfig struct {
	Identifier int
	BotConfs []BotConf
}

type BotConf struct {
	Path string
}

func main() {
	conn := setupConnection(os.Args[1])
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

func setupBot(conf BotConf, portNumber int) *net.TCPConn {
	addrString := "127.0.0.1:" + strconv.Itoa(portNumber)
	fd := []*os.File{os.Stdin,os.Stdout,os.Stderr};
    _, err := os.ForkExec(conf.Path, []string{addrString}, nil, "",fd);
	dieIfError(err, "Error launching bot")
	
	addr, err := net.ResolveTCPAddr(addrString);
	dieIfError(err, "TCP address resolution error")
	conn, err := net.DialTCP("tcp", nil, addr)
	for err != nil && strings.HasSuffix(err.String(), "connection refused") {
		time.Sleep(10000)
		conn, err = net.DialTCP("tcp", nil, addr)
	}
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
