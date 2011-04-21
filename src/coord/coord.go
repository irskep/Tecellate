/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import (
    "agent/link"
    cagent "coord/agent"
    aproxy "coord/agent/proxy"
    "coord/game"
    geo "coord/geometry"
    "coord/config"
    "fmt"
    "logflow"
    "net"
    "netchan"
    "util"
)

type Coordinator struct {
    availableGameState *game.GameState
    peers []*CoordinatorProxy
    rpcSendChannels []chan game.GameStateResponse
    rpcRecvChannels []chan game.GameStateRequest
    conf *config.Config
    Exporter *netchan.Exporter
    listener *net.TCPListener

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
                        rpcRecvChannels: make([]chan game.GameStateRequest, 0),
                        conf: &config.Config{},
                        Exporter: netchan.NewExporter(),
                        rpcRequestsReceivedConfirmation: make(chan int),
                        nextTurnAvailableSignals: make([]chan int, 0),
                        log: logflow.NewSource("coord/?")}
}

func (self *Coordinator) Configure(conf *config.Config) {
    self.conf = conf
    self.conf.Logs.Apply()
    self.availableGameState.Configure(conf)
    self.log = logflow.NewSource(fmt.Sprintf("coord/%d", conf.Identifier))
    self.log.Printf("Configured")
}

func (self *Coordinator) Config() *config.Config {
    return self.conf
}

func (self *Coordinator) Run() {
    // Spawns a bunch of goroutines and exits
    self.StartRPCServer()

    // Run on main thread so we don't need a 'complete' channel
    self.ProcessTurns(nil)
}

func (self *Coordinator) Close() {
    if self.listener != nil {
        self.log.Print("Closing channels")
        self.listener.Close()
        for _, p := range(self.peers) {
            if !closed(p.sendChannel) {
                close(p.sendChannel)
            }
            if !closed(p.recvChannel) {
                close(p.recvChannel)
            }
        }
    }
}

func (self *Coordinator) GetGameState() *game.GameState {
    return self.availableGameState
}

// Set up the server end of an RPC relationship
func (self *Coordinator) AddRPCChannel(newSendChannel chan game.GameStateResponse, newRecvChannel chan game.GameStateRequest) {
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
    newRecvChannel := make(chan game.GameStateRequest)

    // Add a proxy for new peer
    self.peers = append(self.peers, 
                        NewCoordProxy(other.conf.Identifier, self.conf.Identifier, 
                                      other.conf.Address,
                                      newRecvChannel, newSendChannel))

    // Tell peer to listen for RPC requests from me
    other.AddRPCChannel(newSendChannel, newRecvChannel)
}

// REMOTE/PRODUCTION

func (self *Coordinator) NumInitialConns() int {
    // If set up separate-process, Peers will be populated.
    numPeers := len(self.conf.Peers)
    // If set up in-process, rpcSendChannels will be populated.
    // I'M SO SORRY
    if len(self.rpcSendChannels) > numPeers {
        numPeers = len(self.rpcSendChannels)
    }
    return numPeers + len(self.conf.Agents)
}

func (self *Coordinator) RunExporterInitial() {
    go self.RunExporterBlocking(self.NumInitialConns())
}

func (self *Coordinator) RunExporterBlocking(n int) {
    self.log.Print("Listening at ", self.conf.Address, " for " , n)
    
    addr, err := net.ResolveTCPAddr(self.conf.Address)
    if err != nil {
        self.log.Fatal(err)
    }
    
    self.listener, err = net.ListenTCP(addr.Network(), addr)
    if err != nil {
        self.log.Fatal(err)
    }
    // RACE CONDITION!
    for i := 0; i < n; i++ {
        conn, err := self.listener.AcceptTCP()
        self.log.Print("Serving netchan export ", i, " of ", n)
        if err != nil {
            self.log.Fatal("listen:", err)
        }
        conn.SetLinger(0)
        go self.Exporter.ServeConn(conn)
    }
    self.log.Print("Closing listener")
    self.listener.Close()
}

func (self *Coordinator) NewProxy(s *cagent.AgentState) cagent.Agent {
    p2a := make(chan link.Message, 10)
    a2p := make(chan link.Message, 10)

    self.log.Print("Exporting ", fmt.Sprintf("agent_rsp_%d", s.Id))

    err := self.Exporter.Export(fmt.Sprintf("agent_rsp_%d", s.Id), p2a, netchan.Send)
    if err != nil {
        self.log.Fatal(err)
    }

    self.log.Print("Exporting ", fmt.Sprintf("agent_req_%d", s.Id))

    err = self.Exporter.Export(fmt.Sprintf("agent_req_%d", s.Id), a2p, netchan.Recv)
    if err != nil {
        self.log.Fatal(err)
    }

    proxy := aproxy.NewAgentProxy(p2a, a2p)
    proxy.SetState(s)
    return proxy
}

func (self *Coordinator) AddNewProxyFromState(s *cagent.AgentState) {
    self.availableGameState.Agents = append(self.availableGameState.Agents, self.NewProxy(s))
}

func (self *Coordinator) PrepareAgentProxies() {
    for _, ad := range(self.conf.Agents) {
        s := cagent.NewAgentState(ad.Id, 0, *geo.NewPoint(ad.X, ad.Y), cagent.Energy(ad.Energy))
        self.AddNewProxyFromState(s)
    }
}

func (self *Coordinator) PrepareCoordProxies() {
    for _, id := range(self.conf.Peers) {
        self.ExportRemote(id)
    }
}

func (self *Coordinator) ConnectCoordProxies() {
    self.log.Print(self.conf.Peers)
    for address, id := range(self.conf.Peers) {
        self.ConnectToRPCServer(id, address)
    }
}

func (self *Coordinator) ExportRemote(otherID int) {
    ch_recv := make(chan game.GameStateRequest)
    ch_send := make(chan game.GameStateResponse)
    
    err := self.Exporter.Export(fmt.Sprintf("coord_req_%d", otherID), ch_recv, netchan.Recv)
    if err != nil {
	    self.log.Fatal(err)
	}

    err = self.Exporter.Export(fmt.Sprintf("coord_rsp_%d", otherID), ch_send, netchan.Send)
	if err != nil {
	    self.log.Fatal(err)
	}

	self.AddRPCChannel(ch_send, ch_recv)
}

func (self *Coordinator) ConnectToRPCServer(otherID int, otherAddress string) {
    ch_send := make(chan game.GameStateRequest)
    ch_recv := make(chan game.GameStateResponse)

    self.log.Printf("Importing coord_req_%d from %v", otherID, otherAddress)
    
    imp := util.MakeImporterWithRetry("tcp", otherAddress, 10, self.log)

	err := imp.Import(fmt.Sprintf("coord_req_%d", self.conf.Identifier), ch_send, netchan.Send, 1)
	if err != nil {
	    self.log.Fatal(err)
	}

	err = imp.Import(fmt.Sprintf("coord_rsp_%d", self.conf.Identifier), ch_recv, netchan.Recv, 1)
	if err != nil {
	    self.log.Fatal(err)
	}
    
    np := NewCoordProxy(otherID, self.conf.Identifier, self.conf.Address, ch_send, ch_recv)
    self.peers = append(self.peers, np)
}
