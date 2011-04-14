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
    "logflow"
    "net"
    "netchan"
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

func (self CoordinatorSlice) Chain() {
    for i, c := range(self) {
        if i < len(self)-1 {
            logflow.Printf("main", "Connect %d to %d locally", i, i+1)
            c.ConnectToLocal(self[i+1])
        }
        if i > 0 {
            logflow.Printf("main", "Connect %d to %d locally", i, i-1)
            c.ConnectToLocal(self[i-1])
        }
    }
}

func (self CoordinatorSlice) ChainTCP() {
    for i, c := range(self) {
        c.InitTCP()
        if i < len(self)-1 {
            c.ExportRemote(i+1)
        }
        if i > 0 {
            c.ExportRemote(i-1)
        }
        c.ListenForRPCConnections()
    }
    for i, c := range(self) {
        if i < len(self)-1 {
            logflow.Printf("main", "Connect %d to %d over TCP", i, i+1)
            c.ConnectToRPCServer(i+1)
        }
        if i > 0 {
            logflow.Printf("main", "Connect %d to %d over TCP", i, i-1)
            c.ConnectToRPCServer(i-1)
        }
    }
    for _, c := range(self) {
        c.listener.Close()
    }
}

/* Coordinator type */

type Coordinator struct {
    availableGameState *game.GameState
    peers []*CoordinatorProxy
    rpcSendChannels []chan interface{}
    rpcRecvChannels []chan interface{}
    conf *config.Config
    listener net.Listener
    exporter *netchan.Exporter
    
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

    log logflow.Logger
}

/* Initialization */

// Create a new Coordinator. Initialize but do not fill the data structures.
func NewCoordinator() *Coordinator {
    return &Coordinator{availableGameState: game.NewGameState(),
                        peers: make([]*CoordinatorProxy, 0),
                        rpcSendChannels: make([]chan interface{}, 0),
                        rpcRecvChannels: make([]chan interface{}, 0),
                        rpcRequestsReceivedConfirmation: make(chan int),
                        nextTurnAvailableSignals: make([]chan int, 0),
                        log: logflow.NewSource("coord/?")}
}

func (self *Coordinator) Configure(conf *config.Config) {
    self.conf = conf
    self.availableGameState.Configure(conf)
    self.log = logflow.NewSource(fmt.Sprintf("coord/%d", conf.Identifier))
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

// REMOTE/PRODUCTION

func (self *Coordinator) InitTCP() {
    addr_string := fmt.Sprintf("127.0.0.1:%d", 8000+self.conf.Identifier)
    self.log.Println("Listening at", addr_string)
    addr, _ := net.ResolveTCPAddr(addr_string)
    lstn, err := net.ListenTCP(addr.Network(), addr)
    self.listener = lstn
    if err != nil {
        self.log.Fatal(err)
    }
    self.exporter = netchan.NewExporter()
}

func (self *Coordinator) ExportRemote(otherID int) {
    ch_recv := make(chan interface{})
    ch_send := make(chan interface{})
    
    err := self.exporter.Export(fmt.Sprintf("coord_req_%d", otherID), ch_recv, netchan.Recv)
    if err != nil {
	    self.log.Fatal(err)
	}
	
    err = self.exporter.Export(fmt.Sprintf("coord_rsp_%d", otherID), ch_send, netchan.Send)
	if err != nil {
	    self.log.Fatal(err)
	}
	
    self.peers = append(self.peers, NewCoordProxy(otherID, self.conf.Identifier, ch_send, ch_recv))
}

func (self *Coordinator) ListenForRPCConnections() {
    go self.exporter.Serve(self.listener)
}

func (self *Coordinator) ConnectToRPCServer(otherID int) {
    ch_send := make(chan interface{})
    ch_recv := make(chan interface{})
    
    imp, err := netchan.Import("tcp", fmt.Sprintf("127.0.0.1:%d", 8000+otherID))
    if err != nil {
	    self.log.Fatal(err)
	}
	
	err = imp.Import(fmt.Sprintf("coord_req_%d", otherID), ch_send, netchan.Send, 1)
	if err != nil {
	    self.log.Fatal(err)
	}
	
	err = imp.Import(fmt.Sprintf("coord_rsp_%d", otherID), ch_recv, netchan.Recv, 1)
	if err != nil {
	    self.log.Fatal(err)
	}
	
	self.AddRPCChannel(ch_send, ch_recv)
}
