package coord

import (
    "coord/game"
)

func (self *Coordinator) StartRPCServer() {
    for i, requestChannel := range(self.rpcRecvChannels) {
        responseChannel := self.rpcSendChannels[i]
        go self.serveRPCRequestsOnChannels(self.peers[i].Identifier, requestChannel, responseChannel, self.nextTurnAvailableSignals[i])
    }
}

func (self *Coordinator) serveRPCRequestsOnChannels(identifier int,
                                                    requestChannel chan game.GameStateRequest,
                                                    responseChannel chan game.GameStateResponse,
                                                    nextTurnAvailable chan int) {
    for i := 0 ; i < self.conf.MaxTurns; i++ {
        self.log.Printf("(%d) Waiting for turn %d", identifier, i)
        // Wait for turn i to become available
        <- nextTurnAvailable
        
        self.log.Printf("(%d) Now serving a request for turn %d", identifier, i)
        
        // Read a request
        request := <- requestChannel
        
        // Send the response
        responseChannel <- self.availableGameState.MakeRPCResponse(request)
        
        // Send an RPC request confirmation down the pipes so the
        // processing loop knows when it is allowed to proceed
        self.rpcRequestsReceivedConfirmation <- request.Turn
    }
}
