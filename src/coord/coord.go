/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import (
    "coord/game"
    "coord/config"
    "fmt"
    "log"
    "os"
)

/* Coordinator bucket and convenience methods */

type CoordinatorSlice []*Coordinator

func (self CoordinatorSlice) Run() {
    // This channel will receive one 'true' for each process completion
    complete := make(chan bool)
    
    // Start the necessary threads
    for _, c := range(self) {
        c.StartRPCServer()
        go c.ProcessTurns(complete)
    }
    
    // Wait for processing to complete
    for _, _ = range(self) {
        <- complete
    }
}

/* Coordinator type */

type Coordinator struct {
    availableGameState *game.GameState
    peers []*CoordinatorProxy
    rpcSendChannels []chan interface{}
    rpcRecvChannels []chan interface{}
    conf *config.Config
    
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
    
    log *log.Logger
}

/* Initialization */

// Create a new Coordinator. Initialize but do not fill the data structures.
func NewCoordinator() *Coordinator {
    return &Coordinator{availableGameState: game.NewGameState(), 
                        peers: make([]*CoordinatorProxy, 0),
                        rpcSendChannels: make([]chan interface{}, 0),
                        rpcRecvChannels: make([]chan interface{}, 0),
                        conf: nil,
                        rpcRequestsReceivedConfirmation: make(chan int),
                        nextTurnAvailableSignals: make([]chan int, 0),
                        log: nil}
}

func (self *Coordinator) Configure(conf *config.Config) {
    self.conf = conf
    self.availableGameState.Configure(conf)
    self.log = log.New(os.Stdout, fmt.Sprintf("%d: ", conf.Identifier), 0)
    self.log.Printf("Configured")
}

func (self *Coordinator) Run() {
    // Spawns a bunch of goroutines and exits
    self.StartRPCServer()
    
    // Run on main thread so we don't need a 'complete' channel
    self.ProcessTurns(nil)
}

// LOCAL/TESTING

// Set up a connection with another coordinator in the same process.
func (self *Coordinator) ConnectToLocal(other *Coordinator) {
    // We communicate over this channel instead of a netchan
    newSendChannel := make(chan interface{})
    newRecvChannel := make(chan interface{})
    
    // Add a proxy for new peer
    self.peers = append(self.peers, NewCoordProxy(other.conf.Identifier, self.conf.Identifier, newSendChannel, newRecvChannel))
    
    // Tell peer to listen for RPC requests from me
    other.AddRPCChannel(newRecvChannel, newSendChannel)
}

// Set up the server end of an RPC relationship
func (self *Coordinator) AddRPCChannel(newSendChannel chan interface{}, newRecvChannel chan interface{}) {
    // Add the given channel to a list of RPC channels to be read later
    self.rpcSendChannels = append(self.rpcSendChannels, newSendChannel)
    self.rpcRecvChannels = append(self.rpcRecvChannels, newRecvChannel)
    
    // Also add a channel-as-lock to correspond to this RPC channel.
    // Every time a new turn is available, the turn's number is sent down this channel.
    // There is one channel per RPC server, so the processing loop sends k ints to k RPC threads.
    self.nextTurnAvailableSignals = append(self.nextTurnAvailableSignals, make(chan int))
}

// LOCAL THAT MIMICS REMOTE BETTER

type LocalPeeringRequest struct {
    
}

func (self *Coordinator) ListenLocal(k int) {
    go func() {
        for i := 0; i < k; i++ {
            
        }
    }()
}

// REMOTE/PRODUCTION

// TCP-based version of ConnectToLocal.
func (self *Coordinator) ConnectToRemote(address []byte) {
    
}

// Rather than having ConnectToRemote call some specific function, have it dial a
// TCP port which we will listen on, accept connections on, and add those connections
// to self.rpcSendChannels as netchans.
func (self *Coordinator) ListenForRPCConnectionSetupRequests(address []byte) {
    
}
