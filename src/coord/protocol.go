package main

import (
	"fmt"
	"easynet"
	"json"
	"net"
	"time"
)

type CompletionNotification struct {
	Identifier int
	LastProcessedTurn int
}

type Request struct {
	Identifier int
	Turn int
	Command string
}

type BotInfo struct {
	x int
	y int
}

type RespondNodeInfo struct {
	Identifier int
	Turn int
	BotData []BotInfo
}

func listenForMaster(connectionToMaster *net.TCPConn) {
	msg := string(easynet.ReceiveFrom(connectionToMaster))
	if msg == "begin" {
		fmt.Printf("%d is primary\n", config.Identifier)
		primary = true
		broadcastValid()
	}
}

func listenForPeer() {
	fmt.Printf("%d serving requests\n", config.Identifier)
	for data := range(listenServe) {
		r := new(Request)
		err := json.Unmarshal(data, r)
		easynet.DieIfError(err, "JSON error")
		switch {
		case r.Command == "Begin" && primary == false && waitingForStart == true:
			fmt.Printf("%d handle Begin from %d\n", config.Identifier, r.Identifier)
			waitingForStart = false
			go processNodes()
		case r.Command == "GetNodes":
			fmt.Printf("%d handle GetNodes from %d\n", config.Identifier, r.Identifier)
			for respondingToRequestsFor < r.Turn {
				fmt.Printf("%d not ready for GetNodes\n", config.Identifier)
				time.Sleep(100000)
			}
			fmt.Printf("%d ready for GetNodes\n", config.Identifier)
			info := new(RespondNodeInfo)
			info.Identifier = config.Identifier
			info.Turn = respondingToRequestsFor
			info.BotData = nil
			infoString, err := json.Marshal(info)
			easynet.DieIfError(err, "JSON marshal error")
			adjsServe[r.Identifier].Write(infoString)
			fmt.Printf("%d sent GetNodes response to %d\n", config.Identifier, r.Identifier)
		}
	}
}

func processNodes() {
	fmt.Printf("%d processing nodes\n", config.Identifier)
	for i := 0; i < config.NumTurns; i++ {
		respondingToRequestsFor = i
		fmt.Printf("%d turn %d\n", config.Identifier, i)
		for j, conn := range(adjsRequest) {
			fmt.Printf("%d turn %d, request neighbor %d\n", config.Identifier, i, j)
			r := new(Request)
			r.Identifier = config.Identifier
			r.Turn = respondingToRequestsFor
			r.Command = "GetNodes"
			
			rData, err := json.Marshal(r)
			easynet.DieIfError(err, "JSON marshal error")
			conn.Write(rData)
			
			info := new(RespondNodeInfo)
			err = json.Unmarshal(easynet.ReceiveFrom(conn), info)
			easynet.DieIfError(err, "JSON unmarshal error")
		}
	}
	complete <- true
}

func broadcastValid() {
	note := new(Request)
	note.Identifier = config.Identifier
	note.Turn = respondingToRequestsFor
	note.Command = "Begin"
	data, err := json.Marshal(note)
	easynet.DieIfError(err, "JSON marshal error")
	
	for i, conn := range(adjsRequest) {
		fmt.Printf("%d broadcasting to %d\n", config.Identifier, i)
		conn.Write(data)
	}
	waitingForStart = false
	go processNodes()
	
	time.Sleep(10000)
	go listenForPeer()
}
