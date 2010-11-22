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

type coordMap map[int]*net.TCPConn
var adjsServe coordMap
var adjsRequest coordMap
var config *ttypes.CoordConfig

func main() {
	listener := easynet.HostWithAddress(os.Args[1])
	defer listener.Close()
	conn := easynet.Accept(listener)
	defer conn.Close()
	
	config = new(ttypes.CoordConfig)
	err := json.Unmarshal(easynet.ReceiveFrom(conn), config)
	easynet.DieIfError(err, "JSON error")
	fmt.Println(config)
	
	conn.Write([]uint8("connected"))
	
	connections, botsComplete := setupBots(config)
	
	adjsServe = make(coordMap, len(config.AdjacentCoords))
	adjsRequest = make(coordMap, len(config.AdjacentCoords))
	adjsComplete := connectToAdjacents(listener)
	
	<- adjsComplete
	<- botsComplete
	
	for _, c := range(connections) {
		fmt.Printf("%s\n", easynet.ReceiveFrom(c))
	}
	
	conn.Write([]uint8("finished"))
}

func setupBot(conf ttypes.BotConf, portNumber int) *net.TCPConn {
	addrString := "127.0.0.1:" + strconv.Itoa(portNumber)
	fd := []*os.File{os.Stdin,os.Stdout,os.Stderr};
    _, err := os.ForkExec(conf.Path, []string{addrString}, nil, "",fd);
	easynet.DieIfError(err, "Error launching bot")
	
	return easynet.Dial(addrString)
}

func setupBots(config *ttypes.CoordConfig) ([]*net.TCPConn, chan bool) {
	connections := make([]*net.TCPConn, len(config.BotConfs))
	basePort := new(int)
	botComplete := make(chan bool)
	fmt.Sscanf(os.Args[1], "127.0.0.1:%d", basePort)
	go func() {
		for ix, b := range(config.BotConfs) {
			connections[ix] = setupBot(b, *basePort + ix + 1)
		}
		botComplete <- true
	}()
	
	return connections, botComplete
}

func connectToAdjacents(listener *net.TCPListener) chan bool {
	serveFound := make(chan int)
	requestFound := make(chan int)
	allDone := make(chan bool)
	
	go func() {
		for _, adj := range(config.AdjacentCoords) {
			go func() {
				adjsRequest[adj.Identifier] = easynet.Dial(adj.Address)
				adjsRequest[adj.Identifier].Write([]uint8(strconv.Itoa(config.Identifier)))
				requestFound <- adj.Identifier
			}()
			go func() {
				newConn := easynet.Accept(listener)
				identifier, err := strconv.Atoi(string(easynet.ReceiveFrom(newConn)))
				easynet.DieIfError(err, "String conversion error")
				adjsServe[identifier] = newConn
				serveFound <- identifier
			}()
		}
	
		for i := 0; i < len(config.AdjacentCoords); i++ {
			<-requestFound
			<-serveFound
		}
	
		fmt.Printf("%d is connected to all neighbors\n", config.Identifier)
		
		allDone <- true
	}()
	
	return allDone
}
