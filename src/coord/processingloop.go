package coord

import geo "coord/geometry"

// import (
//     "coord/game"
// )

import (
    "rand"
    "time"
)

func (self *Coordinator) ProcessTurns(complete chan bool) {
    for i := 0; i <3 /* <3 <3 <3 */; i++ {  // TODO: THREE TIMES IS ARBITRARY AND FOR TESTING

        self.log.Printf("Making turn %d available", i)
        for pi, _ := range(self.peers) {
            self.nextTurnAvailableSignals[pi] <- i
        }

        responses := self.peerDataForTurn(i)
        _ = self.transformsForNextTurn(responses)

        if (self.conf.RandomlyDelayProcessing) {
            time.Sleep(int64(float64(1e9)*rand.Float64()))
        }

        // Wait for all RPC requests from peers to go through the other goroutine
        for _, _ = range(self.peers) {
            <- self.rpcRequestsReceivedConfirmation
        }

        self.availableGameState.Advance()
        //  i, agent
        self.log.Println(self.availableGameState.Agents)
        for _, prox := range(self.availableGameState.Agents) {
            self.log.Println(prox)
            // agent.Apply(transforms[i])
        }
    }

    self.log.Printf("Sending complete")

    if complete != nil {
        complete <- true
    }
}

func (self *Coordinator) peerDataForTurn(turn int) []*GameStateResponse {
    responses := make([]*GameStateResponse, len(self.peers))
    responsesReceived := make(chan bool)
    for p, _ := range(self.peers) {
        go func(peerIndex int) {
            responses[peerIndex] = self.peers[peerIndex].RequestStatesInBox(turn, geo.Point{0,0}, geo.Point{0,0})
            responsesReceived <- true
        }(p)
    }

    for _, _ = range(self.peers) {
        <- responsesReceived
    }
    return responses
}

type ProspectiveMap *bool   // Make this a struct later

func (self *Coordinator) transformsForNextTurn(peerData []*GameStateResponse) ProspectiveMap {
    agents := self.availableGameState.Agents
    for _, agent := range(agents) {
        success := agent.Turn()
        state := agent.State()
        self.log.Printf("%v, %s", success, state)
    }
    // transforms = make([]*StateTransform, len(agents))
    // Do some magic to make the transforms
    // Return transforms
    return nil;
}
