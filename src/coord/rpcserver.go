package coord

func (self *Coordinator) StartRPCServer() {
    for i, requestChannel := range(self.rpcChannels) {
        go self.serveRPCRequestsOnChannel(requestChannel, self.nextTurnAvailableSignals[i])
    }
}

func (self *Coordinator) serveRPCRequestsOnChannel(requestChannel chan interface{},
                                                   nextTurnAvailable chan int) {
    for i := 0 ; ; i++ {    // Spin forever. Process will exit without our help.
        self.log.Printf("Waiting for turn %d", i)
        // Wait for turn i to become available
        <- nextTurnAvailable
        
        // Read a request
        request := (<- requestChannel).(GameStateRequest)
        
        // Build a response object
        self.log.Printf("Sender %d asked for %d, sending %d", request.SenderIdentifier, request.Turn, i)
        
        // Send the response
        requestChannel <- GameStateResponse{i, nil}
        
        // Send an RPC request confirmation down the pipes so the
        // processing loop knows when it is allowed to proceed
        self.rpcRequestsReceivedConfirmation <- request.Turn
    }
}
