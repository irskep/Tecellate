package main

import (
	"fmt"
	"os"
	"net"
	"strconv"
	"easynet"
	"sync"
	"ttypes"
)

// Configuration
var config *ttypes.CoordConfig

// Connections to neighbors
type CoordMap map[int]*net.TCPConn
var adjsServe CoordMap
var adjsRequest CoordMap
var listenServe chan []uint8

// Which turn we just completed and what the state of everything on that turn was
var respondingToRequestsFor int
var botStates []*BotState

// Global state
var primary bool
var complete chan bool
var processing bool

// Avoid race conditions
var dataLock sync.RWMutex

func main() {
	// Print a nice separator at the end of execution so that 'make fancyrun' looks good
	defer fmt.Println("-------------------")
	
	// Initialize globals
	primary = false
	processing = false
	complete = make(chan bool)
	
	// Begin listening on the port passed on the command line
	listener := easynet.HostWithAddress(os.Args[1])
	defer listener.Close()
	
	// Set up a TCP connection to the master
	connectionToMaster := easynet.Accept(listener)
	defer connectionToMaster.Close()
	
	// Read configuration from the connection to master
	config = new(ttypes.CoordConfig)
	easynet.ReceiveJson(connectionToMaster, config)
	respondingToRequestsFor = 0
	
	// Confirm configuration
	connectionToMaster.Write([]uint8("connected"))
	
	// Set up bots and connect to adjacent coordinators
	setupAll(listener)
	// Remember to terminate all child processes upon termination of this process
	defer killChildren()
	
	// Close all connections when program terminates
	for _, conn := range(adjsServe) {
		defer conn.Close()
	}
	
	// Confirm setup
	connectionToMaster.Write([]uint8("setup complete"))
	
	fmt.Printf("%d sees data at start as: \n%v\n    grid: %v\n", config.Identifier, botInfosForNeighbor(0), config.Terrain)
	
	// Transition to both listening states
	go listenForMaster(connectionToMaster)
	go listenForPeer()
	
	// Wait for those loops to exit
	<-complete
	<-complete
	
	fmt.Printf("%d sees data at end as: \n%v\n    grid: %v\n", config.Identifier, botInfosForNeighbor(0), config.Terrain)
	
	// Confirm termination
	finalTally := new(ttypes.Finish)
	for _, s := range(botStates) {
		if s.Killed == false { finalTally.NumBots += 1 }
	}
	easynet.SendJson(connectionToMaster, finalTally)
}

// Fork a bot and open a TCP connection to it
func setupBot(conf ttypes.BotConf, portNumber int) *net.TCPConn {
	addrString := "127.0.0.1:" + strconv.Itoa(portNumber)
	fd := []*os.File{os.Stdin,os.Stdout,os.Stderr};
    _, err := os.ForkExec(conf.Path, []string{addrString}, nil, "",fd);
	easynet.DieIfError(err, "Error launching bot")
	
	return easynet.Dial(addrString)
}

// Fork and configure all bots
func setupBots() (chan bool) {
	botStates = make([]*BotState, len(config.BotConfs))
	basePort := new(int)
	botComplete := make(chan bool)
	fmt.Sscanf(os.Args[1], "127.0.0.1:%d", basePort)
	go func() {
		for ix, b := range(config.BotConfs) {
			s := new(BotState)
			botStates[ix] = s
			s.Conn = setupBot(b, *basePort + ix + 1)
			s.Info = ttypes.BotInfo{b.X, b.Y, ""}
		}
		botComplete <- true
	}()
	
	return botComplete
}

// Set up incoming and outgoing TCP connections to all adjacent coordinators
func connectToAdjacents(listener *net.TCPListener) chan bool {
	adjsServe = make(CoordMap, len(config.AdjacentCoords))
	adjsRequest = make(CoordMap, len(config.AdjacentCoords))
	listenServe = make(chan []uint8)
	
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
				easynet.TieConnToChannel(newConn, listenServe)
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

// Asynchronously set up neighbors and bots
func setupAll(listener *net.TCPListener) {
	botsComplete := setupBots()
	adjsComplete := connectToAdjacents(listener)
	
	<- adjsComplete
	<- botsComplete
	
	for _, s := range(botStates) {
		fmt.Printf("%s\n", easynet.ReceiveFrom(s.Conn))
	}
}

func killChild(s *BotState) {
	req := new(ttypes.BotMoveRequest)
	req.Kill = true
	easynet.SendJson(s.Conn, req)
}

// Kill all child processes by asking nicely
func killChildren() {
	for _, s := range(botStates) { killChild(s) }
}
