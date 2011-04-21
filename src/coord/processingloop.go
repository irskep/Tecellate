package coord

import game "coord/game"

func (self *Coordinator) ProcessTurns(complete chan bool) {
    self.log.Println("My agents:", self.availableGameState.Agents)
    for i := 0; i < self.conf.MaxTurns; i++ {

        self.log.Printf("Making turn %d available", i)
        for pi, _ := range(self.peers) {
            self.nextTurnAvailableSignals[pi] <- i
        }

        responses := self.peerDataForTurn(i)
        transforms, messages, myMessages := self.transformsForNextTurn(responses)

        // Wait for all RPC requests from peers to go through the other goroutine
        for _, _ = range(self.peers) {
            <- self.rpcRequestsReceivedConfirmation
        }

        self.availableGameState.Advance(transforms, messages, myMessages)
    }

    self.log.Printf("Sending complete")
    
    for _, a := range(self.availableGameState.Agents) {
        a.Stop()
    }

    if complete != nil {
        complete <- true
    }
}

func (self *Coordinator) peerDataForTurn(turn int) []*game.GameStateResponse {
    responses := make([]*game.GameStateResponse, len(self.peers))
    responsesReceived := make(chan bool)
    for p, _ := range(self.peers) {
        go func(peerIndex int) {
            responses[peerIndex] = self.peers[peerIndex].RequestStatesInBox(turn, *self.conf.BottomLeft, *self.conf.TopRight)
            responsesReceived <- true
        }(p)
    }

    for _, _ = range(self.peers) {
        <- responsesReceived
    }
    return responses
}
