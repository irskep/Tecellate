package coord

func (self *Coordinator) StartRPCServer() {
    for i, requestChannel := range(self.rpcRecvChannels) {
        responseChannel := self.rpcSendChannels[i]
        go self.serveRPCRequestsOnChannels(self.peers[i].Identifier, requestChannel, responseChannel, self.nextTurnAvailableSignals[i])
    }
}

func (self *Coordinator) serveRPCRequestsOnChannels(identifier int,
                                                    requestChannel chan interface{},
                                                    responseChannel chan interface{},
                                                    nextTurnAvailable chan int) {
    for i := 0 ; ; i++ {    // Spin forever. Process will exit without our help.
        self.log.Printf("(%d) Waiting for turn %d", identifier, i)
        // Wait for turn i to become available
        <- nextTurnAvailable
        
        self.log.Printf("(%d) Now serving a request for turn %d", identifier, i)
        
        // Read a request
        request := (<- requestChannel).(GameStateRequest)
        
        // Build a response object
        self.log.Printf("Sender %d asked for %d, sending %d", request.SenderIdentifier, self.availableGameState.Turn, i)
        
        // Send the response
        responseChannel <- GameStateResponse{i, nil}
        
        // Send an RPC request confirmation down the pipes so the
        // processing loop knows when it is allowed to proceed
        self.rpcRequestsReceivedConfirmation <- request.Turn
    }
}
