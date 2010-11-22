package main

import (
	"fmt"
	"json"
	"os"
	"net"
	"strconv"
	"../easynet"
	"../ttypes"
)

func main() {
	conn := easynet.HostWithAddress(os.Args[1])
	defer conn.Close()
	
	config := new(ttypes.CoordConfig)
	err := json.Unmarshal(easynet.ReceiveFrom(conn), config)
	easynet.DieIfError(err, "JSON error")
	fmt.Println(config.Terrain)
	
	connections := setupBots(config)
	
	for _, c := range(connections) {
		fmt.Printf("%s\n", easynet.ReceiveFrom(c))
	}
	
	conn.Write([]uint8("ok"))
}

func setupBot(conf ttypes.BotConf, portNumber int) *net.TCPConn {
	addrString := "127.0.0.1:" + strconv.Itoa(portNumber)
	fd := []*os.File{os.Stdin,os.Stdout,os.Stderr};
    _, err := os.ForkExec(conf.Path, []string{addrString}, nil, "",fd);
	easynet.DieIfError(err, "Error launching bot")
	
	return easynet.Dial(addrString)
}

func setupBots(config *ttypes.CoordConfig) []*net.TCPConn {
	connections := make([]*net.TCPConn, len(config.BotConfs))
	basePort := new(int)
	fmt.Sscanf(os.Args[1], "127.0.0.1:%d", basePort)
	for ix, b := range(config.BotConfs) {
		connections[ix] = setupBot(b, *basePort + ix + 1)
	}
	return connections
}
