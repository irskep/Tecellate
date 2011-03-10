package coord

import geo "coord/geometry"

func (self *Coordinator) ProcessTurns(complete chan bool) {
    for i := 0; i <3 /* <3 <3 <3 */; i++ {  // TODO: THREE TIMES IS ARBITRARY AND FOR TESTING
        
        // Signal the availability of turn i to the RPC servers
        for pi, _ := range(self.peers) {
            self.nextTurnAvailableSignals[pi] <- i
        }
        
        for _, peer := range(self.peers) {
            // Probably actually don't want this to be blocking...
            // Also, STORE THE RESULT AND DO SOMETHING WITH IT.
            _ = peer.RequestStatesInBox(i, geo.Point{0,0}, geo.Point{0,0})
        }
        
        // Process new data
        // BLAH BLAH BLAH BLAH BLAH
        
        // Wait for all RPC requests from peers to go through the other goroutine
        for _, _ = range(self.peers) {
            <- self.rpcRequestsReceivedConfirmation
        }
    }
    complete <- true
}
