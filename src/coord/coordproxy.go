package coord

import geo "coord/geometry"
import game "coord/game"

import (
    "fmt"
    "logflow"
    "time"
)

var timeout int64 = 5*1e9

type GameStateRequest struct {
    SenderIdentifier int
    Turn int
    BottomLeft geo.Point
    TopRight geo.Point
}

type CoordinatorProxy struct {
    Identifier int
    parentIdentifier int
    sendChannel chan interface{}
    recvChannel chan interface{}
    log logflow.Logger
}

func NewCoordProxy(identifier int, parentIdentifier int, sendChan chan interface{}, recvChan chan interface{}) *CoordinatorProxy {
    return &CoordinatorProxy{identifier, 
                             parentIdentifier, 
                             sendChan, 
                             recvChan, 
                             logflow.NewSource(fmt.Sprintf("coordproxy/%d/%d: ", parentIdentifier, identifier))}
}

func (self *CoordinatorProxy) request(request interface{}) interface{} {
    self.sendChannel <- request
    
    timeout := time.NewTicker(timeout)
    select {
    case response := <- self.recvChannel:
        timeout.Stop()
        return response
    case <-timeout.C:
        timeout.Stop()
    }
    panic("RPC request timeout")
}

func (self *CoordinatorProxy) RequestStatesInBox(turn int,
                                                 bottomLeft geo.Point,
                                                 topRight geo.Point) *game.GameStateResponse {
    request := GameStateRequest{self.parentIdentifier, turn, bottomLeft, topRight}
    self.log.Printf("req: %v", request)
    response := self.request(request)
    self.log.Printf("rsp: %v", response)
    return (response.(game.GameStateResponse)).CopyToHeap()
}

func (self *CoordinatorProxy) SendComplete() {
}
