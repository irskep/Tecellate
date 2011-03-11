package coord

import geo "coord/geometry"

import (
    "coord/game"
)

import (
    "rand"
    "time"
)

func (self *Coordinator) ProcessTurns(complete chan bool) {
    for i := 0; i <3 /* <3 <3 <3 */; i++ {  // TODO: THREE TIMES IS ARBITRARY AND FOR TESTING
        
        self.log.Printf("Making turn %d available", i)
        
        // Signal the availability of turn i to the RPC servers
        for pi, _ := range(self.peers) {
            self.nextTurnAvailableSignals[pi] <- i
        }
        
        responses := make([]*GameStateResponse, len(self.peers))
        responsesReceived := make(chan bool)
        for p, _ := range(self.peers) {
            go func(turn int, peerIndex int) {
                responses[peerIndex] = self.peers[peerIndex].RequestStatesInBox(turn, geo.Point{0,0}, geo.Point{0,0})
                responsesReceived <- true
            }(i, p)
        }
        
        for _, _ = range(self.peers) {
            <- responsesReceived
        }
        
        // Process new data
        nextState := self.nextGameState(responses)
        if (self.conf.RandomlyDelayProcessing) {
            time.Sleep(int64(float64(1e9)*rand.Float64()))
        }
        
        // Wait for all RPC requests from peers to go through the other goroutine
        for _, _ = range(self.peers) {
            <- self.rpcRequestsReceivedConfirmation
        }
        
        self.applyState(nextState)
    }
    
    self.log.Printf("Sending complete")
    
    if complete != nil {
        complete <- true
    }
}

type ProspectiveMap *bool   // Make this a struct later

func (self *Coordinator) buildProspectiveMap(peerData []*GameStateResponse) ProspectiveMap {
    agents := self.availableGameState.Agents
    for _, agent := range(agents) {
        success := agent.Turn()
        self.log.Printf("%u", success)
    }
    return nil;
}

func (self *Coordinator) nextGameState(peerData []*GameStateResponse) *game.GameState {
    _ = self.buildProspectiveMap(peerData)
    return self.availableGameState.CopyAndAdvance()
}

func (self *Coordinator) applyState(nextState *game.GameState) {
    
}
