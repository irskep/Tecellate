package main

import (
	"fmt"
	"easynet"
	"json"
	"net"
	"time"
	"ttypes"
)

func listenForMaster(connectionToMaster *net.TCPConn) {
	msg, err := easynet.ReceiveFromWithError(connectionToMaster)
	if err != nil {
		fmt.Printf("%d got an error on the connection to master: %v\n", config.Identifier, err)
	} else {
		if string(msg) == "begin" {
			fmt.Printf("%d is the Chosen One!\n", config.Identifier)
			processing = true
			go processNodes()
		}
	}
}

var completionsRemaining int

func listenForPeer() {
	fmt.Printf("%d serving requests\n", config.Identifier)
	completionsRemaining = len(adjsServe)
	for data := range(listenServe) {
		fmt.Println(string(data))
		//Sometimes requests will be stuck together. Here I am separating them.
		//A crappy and hopefully temporary fix.
		splitPoint := 0
		if data[0] == "{"[0] {
			for i := 1; i < len(data); i++ {
				if data[i-1] == "}"[0] && data[i] == "{"[0] {
					splitPoint = i
					break
				}
			}
		} else {
			for i := 1; i < len(data); i++ {
				if data[i] == "{"[0] {
					splitPoint = i
					break
				}
			}
		}
		if splitPoint == 0 {
			handleRequest(data)
		} else {
			fmt.Println("case 2")
			handleRequest(data[0:splitPoint])
			handleRequest(data[splitPoint:len(data)])
		}
	}
}

func handleRequest(data []uint8) {
	if processing == false {
		// The game's afoot!
		processing = true
		go processNodes()
	}
	r := new(Request)
	err := json.Unmarshal(data, r)
	easynet.DieIfError(err, "JSON error")
	switch {
	case r.Command == "GetNodes":
		fmt.Printf("%d handle GetNodes from %d\n", config.Identifier, r.Identifier)
		for respondingToRequestsFor < r.Turn {
			fmt.Printf("%d not ready for GetNodes\n", config.Identifier)
			time.Sleep(10000000)
		}
		dataLock.RLock()
		fmt.Printf("%d ready for GetNodes\n", config.Identifier)
		info := new(RespondNodeInfo)
		info.Identifier = config.Identifier
		info.Turn = respondingToRequestsFor
		info.BotData = botInfosForNeighbor(r.Identifier)
		infoString, err := json.Marshal(info)
		easynet.DieIfError(err, "JSON marshal error")
		adjsServe[r.Identifier].Write(infoString)
		fmt.Printf("%d sent GetNodes response to %d\n", config.Identifier, r.Identifier)
		dataLock.RUnlock()
	case r.Command == "Complete":
		completionsRemaining -= 1
		if completionsRemaining == 0 {
			fmt.Println("Right ho, we're finished here chaps")
			complete<-true
		}
	}
}

func processNodes() {
	fmt.Printf("%d processing nodes\n", config.Identifier)
	
	dataLock.Lock()
	defer dataLock.Unlock()
	for i := 0; i < config.NumTurns; i++ {
		respondingToRequestsFor = i
		dataLock.Unlock()
		
		fmt.Printf("%d starting turn %d\n", config.Identifier, i)
		
		otherInfos := make([]ttypes.BotInfo, len(botStates), len(botStates)*len(adjsServe))
		
		//Copy all infos from botStates into otherInfos
		for i, s := range(botStates) {
			otherInfos[i] = s.Info
		}
		
		//Get updates from neighbors
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
			
			otherInfos = append(otherInfos, info.BotData...)
		}
		
		declareDeaths(otherInfos)
		
		moveBots(otherInfos)
		
		dataLock.Lock()
		//Copy new data back into botStates.
		for i, _ := range(botStates) {
			botStates[i].Info = otherInfos[i]
		}
	}
	broadcastComplete()
	complete<-true
}

func broadcastComplete() {
	note := new(Request)
	note.Identifier = config.Identifier
	note.Turn = respondingToRequestsFor
	note.Command = "Complete"
	
	for i, conn := range(adjsRequest) {
		fmt.Printf("%d broadcasting complete to %d\n", config.Identifier, i)
		easynet.SendJson(conn, note)
	}
}
