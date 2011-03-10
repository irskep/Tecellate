/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import geo "coord/geometry"

import (
    "coord/game"
    "json"
    "log"
)

type Coordinator struct {
    availableGameState *game.GameState
    peers []*CoordinatorProxy
    rpcChannels []chan []byte
    
    // RPC server threads send an ints down this channel representing
    // a turn info request served.
    // So when len(peers) ints are received, the processing loop
    // may continue. (None of this code is written yet.)
    rpcRequestsReceivedConfirmation chan int
    
    // RPC servers block on their corresponding channels
    // to wait for the next turn to be processed.
    // Needed so that when A has not completed and B has, and
    // B requests new data from A, A's RPC server does not provide
    // the old data by mistake.
    nextTurnAvailableSignals []chan int
}

/* Initialization */

// Create a new Coordinator. Initialize but do not fill the data structures.
func NewCoordinator() *Coordinator {
    return &Coordinator{game.NewGameState(), 
                        make([]*CoordinatorProxy, 0),
                        make([]chan []byte, 0),
                        make(chan int),
                        make([]chan int, 0)}
}

// LOCAL/TESTING

// Set up a connection with another coordinator in the same process.
func (self *Coordinator) ConnectToLocal(other *Coordinator) {
    // We communicate over this channel instead of a netchan
    newChannel := make(chan []byte)
    
    // Add a proxy for new peer
    self.peers = append(self.peers, NewCoordProxyWithChannel(newChannel))
    
    // Tell peer to listen for RPC requests from me
    other.AddRPCChannel(newChannel)
}

// Set up the server end of an RPC relationship
func (self *Coordinator) AddRPCChannel(newChannel chan []byte) {
    // Add the given channel to a list of RPC channels to be read later
    self.rpcChannels = append(self.rpcChannels, newChannel)
    
    // Also add a channel-as-lock to correspond to this RPC channel.
    // Every time a new turn is available, the turn's number is sent down this channel.
    // There is one channel per RPC server, so the processing loop sends k ints to k RPC threads.
    self.nextTurnAvailableSignals = append(self.nextTurnAvailableSignals, make(chan int))
}

// REMOTE/PRODUCTION

// TCP-based version of ConnectToLocal.
func (self *Coordinator) ConnectToRemote(address []byte) {
    
}

// Rather than having ConnectToRemote call some specific function, have it dial a
// TCP port which we will listen on, accept connections on, and add those connections
// to self.rpcChannels as netchans.
func (self *Coordinator) ListenForRPCConnectionSetupRequests(address []byte) {
    
}

/* RPC Server */

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
        requestBytes := <- requestChannel
        var request GameStateRequest
        _ = json.Unmarshal(requestBytes, &request)
        
        // Build a response object
        log.Printf("Asked for %d, sending %d", request.Turn, i)
        
        // Send the response
        responseBytes, _ :=  json.Marshal(GameStateResponse{i, nil})
        requestChannel <- responseBytes
        
        // Send an RPC request confirmation down the pipes so the
        // processing loop knows when it is allowed to proceed
        self.rpcRequestsReceivedConfirmation <- request.Turn
    }
}

/* Processing */

func (self *Coordinator) ProcessTurns(complete chan bool) {
    for i := 0; i <3 /* <3 <3 <3 */; i++ {  // TODO: THREE TIMES IS ARBITRARY AND FOR TESTING
        
        // Signal the availability of turn i to the RPC servers
        for pi, _ := range(self.peers) {
            self.nextTurnAvailableSignals[pi] <- i
        }
        
        for _, peer := range(self.peers) {
            // Probably actually don't want this to be blocking...
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
