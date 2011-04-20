package coord

import geo "coord/geometry"
import "coord/game"

import (
    "fmt"
    "logflow"
    "time"
)

var timeout int64 = 5*1e9

type CoordinatorProxy struct {
    Identifier int
    parentIdentifier int
    address string
    sendChannel chan game.GameStateRequest
    recvChannel chan game.GameStateResponse
    log logflow.Logger
}

func NewCoordProxy(identifier int, parentIdentifier int, 
                   address string,
                   sendChan chan game.GameStateRequest, 
                   recvChan chan game.GameStateResponse) *CoordinatorProxy {
    return &CoordinatorProxy{identifier, 
                             parentIdentifier, 
                             address,
                             sendChan, 
                             recvChan, 
                             logflow.NewSource(fmt.Sprintf("coordproxy/%d/%d: ", parentIdentifier, identifier))}
}

func (self *CoordinatorProxy) request(request game.GameStateRequest) game.GameStateResponse {
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
    request := game.GameStateRequest{self.parentIdentifier, turn, bottomLeft, topRight}
    self.log.Printf("req: %v", request)
    response := self.request(request)
    self.log.Printf("rsp: %v", response)
    return response.CopyToHeap()
}
