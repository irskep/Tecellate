package coord

import geo "coord/geometry"

import (
    "log"
    "time"
)

var timeout int64 = 5*1e9

type CoordinatorProxy struct {
    parentIdentifier int
    sendChannel chan interface{}
    recvChannel chan interface{}
}

func NewCoordProxy(parentIdentifier int, sendChan chan interface{}, recvChan chan interface{}) *CoordinatorProxy {
    return &CoordinatorProxy{parentIdentifier, sendChan, recvChan}
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
                                                 topRight geo.Point) *GameStateResponse {
    request := GameStateRequest{self.parentIdentifier, turn, bottomLeft, topRight}
    log.Printf("req: %v", request)
    response := self.request(request)
    log.Printf("rsp: %v", response)
    return (response.(GameStateResponse)).CopyToHeap()
}

func (self *CoordinatorProxy) SendComplete() {
}
