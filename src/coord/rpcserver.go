package coord

import (
    "log"
)

func (self *Coordinator) StartRPCServer() {
    for i, requestChannel := range(self.rpcChannels) {
        go self.serveRPCRequestsOnChannel(requestChannel, self.nextTurnAvailableSignals[i])
    }
}

func (self *Coordinator) serveRPCRequestsOnChannel(requestChannel chan []byte,
                                                   nextTurnAvailable chan int) {
    for i := 0 ; ; i++ {    // Spin forever. Process will exit without our help.
        
        // Wait for turn i to become available
        <- nextTurnAvailable
        
        // Read a request
        request := GameStateRequestFromJson(<- requestChannel)
        
        // Build a response object
        log.Printf("Asked for %d, sending %d", request.Turn, i)
        
        // Send the response
        requestChannel <- GameStateResponseJson(i, nil)
        
        // Send an RPC request confirmation down the pipes so the
        // processing loop knows when it is allowed to proceed
        self.rpcRequestsReceivedConfirmation <- request.Turn
    }
}
