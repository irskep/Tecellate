/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import (
    "coord/game"
    "json"
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
}

/* Initialization */

func NewCoordinator() *Coordinator {
    initialState := game.NewGameState()
    return &Coordinator{initialState, 
                        make([]*CoordinatorProxy, 0),
                        make([]chan []byte, 0),
                        make(chan int)}
}

func (self *Coordinator) ConnectToLocal(other *Coordinator) {
    newChannel := make(chan []byte)
    self.peers = append(self.peers, NewCoordProxyWithChannel(newChannel))
    other.AddRPCChannel(newChannel)
}

func (self *Coordinator) AddRPCChannel(newChannel chan []byte) {
    self.rpcChannels = append(self.rpcChannels, newChannel)
}

/* Running */

func (self *Coordinator) StartRPCServer() {
    for _, channel := range(self.rpcChannels) {
        go self.serveRPCRequestsOnChannel(channel)
    }
}

func (self *Coordinator) serveRPCRequestsOnChannel(channel chan []byte) {
    for i := 0; i <3 /* <3 <3 <3 */; i++ {
        requestBytes := <- channel
        var request GameStateRequest
        _ = json.Unmarshal(requestBytes, &request)
        
        responseBytes, _ :=  json.Marshal(GameStateResponse{i, nil})
        
        channel <- responseBytes
        
        self.rpcRequestsReceivedConfirmation <- request.Turn
    }
}
