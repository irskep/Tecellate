/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import (
    agent "agent"
    cproxy "coord/agent/proxy"
    "coord/game"
    "coord/config"
    "fmt"
    "logflow"
    "net"
    "netchan"
    "os"
    "time"
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
    logflow.Println("main", "Exporting channels")
    ready := make(chan bool)
    for i, c := range(self) {
        if i < len(self)-1 {
            c.ExportRemote(i+1)
        }
        if i > 0 {
            c.ExportRemote(i-1)
        }
        c.RunExporter(ready, len(c.rpcSendChannels)*2+len(c.conf.Agents))
    }
    for _, _ = range(self) {
        <- ready
    }
    logflow.Println("main", "Connecting coordinators")
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
}

func (self CoordinatorSlice) ConnectToLocalAgents(agents map[uint]agent.Agent) {
    for _, c := range(self) {
        for _, ad := range(c.conf.Agents) {
            p := cproxy.RunAgentLocal(agents[ad.Id], ad.X, ad.Y)
            c.availableGameState.Agents = append(c.availableGameState.Agents, p)
        }
    }
}

/* Coordinator type */

type Coordinator struct {
    availableGameState *game.GameState
    peers []*CoordinatorProxy
    rpcSendChannels []chan game.GameStateResponse
    rpcRecvChannels []chan GameStateRequest
    conf *config.Config
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
                        rpcSendChannels: make([]chan game.GameStateResponse, 0),
                        rpcRecvChannels: make([]chan GameStateRequest, 0),
                        exporter: netchan.NewExporter(),
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

func (self *Coordinator) GetGameState() *game.GameState {
    return self.availableGameState
}

// Set up the server end of an RPC relationship
func (self *Coordinator) AddRPCChannel(newSendChannel chan game.GameStateResponse, newRecvChannel chan GameStateRequest) {
    // Add the given channel to a list of RPC channels to be read later
    self.rpcSendChannels = append(self.rpcSendChannels, newSendChannel)
    self.rpcRecvChannels = append(self.rpcRecvChannels, newRecvChannel)

    // Also add a channel-as-lock to correspond to this RPC channel.
    // Every time a new turn is available, the turn's number is sent down this channel.
    // There is one channel per RPC server, so the processing loop sends k ints to k RPC threads.
    self.nextTurnAvailableSignals = append(self.nextTurnAvailableSignals, make(chan int))
}

// LOCAL/TESTING

// Set up a connection with another coordinator in the same process.
func (self *Coordinator) ConnectToLocal(other *Coordinator) {
    // We communicate over this channel instead of a netchan
    newSendChannel := make(chan game.GameStateResponse)
    newRecvChannel := make(chan GameStateRequest)

    // Add a proxy for new peer
    self.peers = append(self.peers, NewCoordProxy(other.conf.Identifier, self.conf.Identifier, newRecvChannel, newSendChannel))

    // Tell peer to listen for RPC requests from me
    other.AddRPCChannel(newSendChannel, newRecvChannel)
}

// REMOTE/PRODUCTION

func (self *Coordinator) ExportRemote(otherID int) {
    ch_recv := make(chan GameStateRequest)
    ch_send := make(chan game.GameStateResponse)

    err := self.exporter.Export(fmt.Sprintf("coord_req_%d", otherID), ch_recv, netchan.Recv)
    if err != nil {
	    self.log.Fatal(err)
	}

    err = self.exporter.Export(fmt.Sprintf("coord_rsp_%d", otherID), ch_send, netchan.Send)
	if err != nil {
	    self.log.Fatal(err)
	}

	self.AddRPCChannel(ch_send, ch_recv)
}

func (self *Coordinator) RunExporter(ready chan bool, numChans int) {
    go func() {
        addr_string := fmt.Sprintf("127.0.0.1:%d", 8000+self.conf.Identifier)
        self.log.Println("Listening at", addr_string)
        addr, _ := net.ResolveTCPAddr(addr_string)
        lstn, err := net.ListenTCP(addr.Network(), addr)
        if err != nil {
            self.log.Fatal(err)
        }
        // There is a race condition here. There is a very slim chance that the
        // main thread will unblock (it is waiting for ready) and yet the call to
        // lstn.Accept() will not have been executed yet, which will cause the
        // client's netchan import to fail.
        // However, the chance is extremely slim.
        ready <- true
        for i := 0; i < numChans; i++ {
        	conn, err := lstn.Accept()
        	if err != nil {
        		self.log.Fatal("listen:", err)
        	}
        	go self.exporter.ServeConn(conn)
        }
    }()
}

func (self *Coordinator) makeImporterWithRetry(network string, remoteaddr string) *netchan.Importer {
    // This method is actually entirely futile because the race condition we're trying
    // to account for happens between listener creation and exporter.ServeConn().
    // An error is only thrown if the listener does not exist, but we must already
    // have a listener to call ServeConn().
    // To really fix this, you have to try sending a message down the pipe and see
    // if it panics.
    var err os.Error
    for i := 0; i < 3; i++ {
        conn, err := net.Dial(network, "", remoteaddr)
        if err == nil {
            return netchan.NewImporter(conn)
        }
        self.log.Print("Netchan import failed, retrying")
        time.Sleep(1e9/2)
    }
    self.log.Print("Netchan import failed three times. Bailing out.")
    self.log.Fatal(err)
    return nil
}

func (self *Coordinator) ConnectToRPCServer(otherID int) {
    ch_send := make(chan GameStateRequest)
    ch_recv := make(chan game.GameStateResponse)

    imp := self.makeImporterWithRetry("tcp", fmt.Sprintf("127.0.0.1:%d", 8000+otherID))

	err := imp.Import(fmt.Sprintf("coord_req_%d", self.conf.Identifier), ch_send, netchan.Send, 1)
	if err != nil {
	    self.log.Fatal(err)
	}

	err = imp.Import(fmt.Sprintf("coord_rsp_%d", self.conf.Identifier), ch_recv, netchan.Recv, 1)
	if err != nil {
	    self.log.Fatal(err)
	}

    self.peers = append(self.peers, NewCoordProxy(otherID, self.conf.Identifier, ch_send, ch_recv))
}
