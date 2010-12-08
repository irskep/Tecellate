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
	
	listenForMoveRequests(conn)
}

func listenForMoveRequests(conn *net.TCPConn) {
	listenServe := make(chan []uint8)
	easynet.TieConnToChannel(conn, listenServe)
	for data := range(listenServe) {
		r := new(ttypes.BotMoveRequest)
		err := json.Unmarshal(data, r)		
		easynet.DieIfError(err, "JSON error")
		
		if r.Kill {
			fmt.Printf("Bot on %s received kill signal\n", os.Args[0])
			os.Exit(0)
		}
		
		fmt.Printf("Bot at %d, %d received messages: %v\n", r.YourX, r.YourY, r.Messages)
		//Do something
		
		response := new(ttypes.BotMoveResponse)
		if r.YourY < 5 {
			response.MoveDirection = "up"
		} else if r.YourY > 5 {
			response.MoveDirection = "down"
		}
		response.BroadcastMessage = fmt.Sprintf("I am %v at %d, %d", os.Args[0], r.YourX, r.YourY)
		responseString, err := json.Marshal(response)
		easynet.DieIfError(err, "JSON marshal error")
		conn.Write(responseString)
	}
}