package main

import (
	"fmt"
	"json"
	"net"
	"os"
	"easynet"
	"ttypes"
)

func main() {
	fmt.Printf("testbot launched on %s\n", os.Args[0])
	listener := easynet.HostWithAddress(os.Args[0])
	defer listener.Close()
	conn := easynet.Accept(listener)
	defer conn.Close()
	
	conn.Write([]uint8("bot setup complete " + os.Args[0]))
	
	go listenForMoveRequests(conn)
}

func listenForMoveRequests(conn *net.TCPConn) {
	listenServe := make(chan []uint8)
	easynet.TieConnToChannel(conn, listenServe)
	for data := range(listenServe) {
		r := new(ttypes.BotMoveRequest)
		err := json.Unmarshal(data, r)
		easynet.DieIfError(err, "JSON error")
		response := new(ttypes.BotMoveResponse)
		response.MoveDirection = "left"
		responseString, err := json.Marshal(response)
		easynet.DieIfError(err, "JSON marshal error")
		conn.Write(responseString)
	}
}