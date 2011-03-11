package coord

import geo "coord/geometry"

import (
    "fmt"
    "log"
    "os"
    "time"
)

var timeout int64 = 5*1e9

type CoordinatorProxy struct {
    Identifier int
    parentIdentifier int
    sendChannel chan interface{}
    recvChannel chan interface{}
    log *log.Logger
}

func NewCoordProxy(identifier int, parentIdentifier int, sendChan chan interface{}, recvChan chan interface{}) *CoordinatorProxy {
    return &CoordinatorProxy{identifier, 
                             parentIdentifier, 
                             sendChan, 
                             recvChan, 
                             log.New(os.Stdout, 
                                     fmt.Sprintf("%d-%d: ", parentIdentifier, identifier),
                                     0)}
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
    self.log.Printf("req: %v", request)
    response := self.request(request)
    self.log.Printf("rsp: %v", response)
    return (response.(GameStateResponse)).CopyToHeap()
}

func (self *CoordinatorProxy) SendComplete() {
}
