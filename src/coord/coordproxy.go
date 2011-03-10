package coord

import geo "coord/geometry"

import (
    "time"
)

type CoordinatorProxy struct {
    parentIdentifier int
    conn chan interface{}
}

func NewCoordProxy(parentIdentifier int, channel chan interface{}) *CoordinatorProxy {
    return &CoordinatorProxy{parentIdentifier, channel}
}

func (self *CoordinatorProxy) RequestStatesInBox(turn int,
                                                 bottomLeft geo.Point,
                                                 topRight geo.Point) GameStateResponse {
    
    self.conn <- GameStateRequest{self.parentIdentifier, turn, bottomLeft, topRight}
    
    timeout := time.NewTicker(5*1e9)
    select {
    case response := <-self.conn:
        timeout.Stop()
        return response.(GameStateResponse)
    case <-timeout.C:
        timeout.Stop()
    }
    panic("RPC request timeout")
}

func (self *CoordinatorProxy) SendComplete() {
}
