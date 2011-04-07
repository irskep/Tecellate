package coord

import game "coord/game"
import geo "coord/geometry"
import cagent "coord/agent"

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
        transforms := self.transformsForNextTurn(responses)
        
        // Stress test to discover race conditions
        if (self.conf.RandomlyDelayProcessing) {
            time.Sleep(int64(float64(1e9)*rand.Float64()))
        }

        // Wait for all RPC requests from peers to go through the other goroutine
        for _, _ = range(self.peers) {
            <- self.rpcRequestsReceivedConfirmation
        }

        self.availableGameState.Advance(transforms)
    }

    self.log.Printf("Sending complete")

    if complete != nil {
        complete <- true
    }
}

func (self *Coordinator) peerDataForTurn(turn int) []*game.GameStateResponse {
    responses := make([]*game.GameStateResponse, len(self.peers))
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

func (self *Coordinator) transformsForNextTurn(peerData []*game.GameStateResponse) []cagent.Transform {
    self.log.Printf("From my neighbors, I see:")
    for _, s := range peerData {
        self.log.Printf("%v", *s)
    }
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))

    for ix, agent := range(agents) {
        _ = agent.Turn()
        state := agent.State()
        t := transformFromState(state)
        t.turn = self.availableGameState.Turn+1
        transforms[ix] = t
    }
    return transforms;
}
