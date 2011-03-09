/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import "coord/game"

type Coordinator struct {
    AvailableGameState *game.GameState
    Peers []*CoordinatorProxy
    RPCChannels []chan []byte
}

func NewCoordinator() *Coordinator {
    initialState := game.NewGameState()
    return &Coordinator{initialState, 
                        make([]*CoordinatorProxy, 0),
                        make([]chan []byte, 0)}
}

func (self *Coordinator) ConnectToLocal(other *Coordinator) {
    newChannel := make(chan []byte)
    self.Peers = append(self.Peers, NewCoordProxyWithChannel(newChannel))
    other.AddRPCChannel(newChannel)
}

func (self *Coordinator) AddRPCChannel(newChannel chan []byte) {
    self.RPCChannels = append(self.RPCChannels, newChannel)
}
