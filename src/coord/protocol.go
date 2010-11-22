package main

import (
	"fmt"
	"easynet"
	"json"
	"net"
)

type TurnFlip struct {
	TurnsRemaining int
}

func listenForStart(connectionToMaster *net.TCPConn) {
	easynet.ReceiveFrom(connectionToMaster)
	//Check message == "begin"
	primary = true
	flipTurns()
}

func listenForTurnFlips(complete chan bool) {
	flip := new(TurnFlip)
	fmt.Printf("%d listening for flips\n", config.Identifier)
	for turnsRemaining > 0 {
		for _, flipTurnConfirmationConn := range(adjsRequest) {
			fmt.Printf(" %d subflip\n", config.Identifier)
			flip.TurnsRemaining = turnsRemaining+1
			for flip.TurnsRemaining >= turnsRemaining {
				flipJson := easynet.ReceiveFrom(flipTurnConfirmationConn)
				err := json.Unmarshal(flipJson, flip)
				easynet.DieIfError(err, "JSON Unmarshal error")
				fmt.Printf("%d received flip notification for turn %d\n", 
						   config.Identifier, flip.TurnsRemaining)
			}
			fmt.Printf(" %d confirm subflip\n", config.Identifier)
		}
		fmt.Printf("%d flipping!\n", config.Identifier)
		processBots()
		flipTurns()
	}
	complete <- true
}

func listenForInfoRequests() {
	
}

func processBots() {
	//Request data from adjsRequest
	//Calculate bot stuff
}

func flipTurns() {
	turnsRemaining -= 1
	flip := new(TurnFlip)
	flip.TurnsRemaining = turnsRemaining
	
	data, err := json.Marshal(flip)
	easynet.DieIfError(err, "JSON Marshal error")
	
	fmt.Printf("%d flipping to turn %d\n", config.Identifier, turnsRemaining)
	for _, flipTurnConfirmationConn := range(adjsServe) {
		flipTurnConfirmationConn.Write(data)
	}
}
