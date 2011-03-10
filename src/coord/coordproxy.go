package coord

import geo "coord/geometry"

import (
    "time"
)

type CoordinatorProxy struct {
    parentIdentifier int
    conn chan []byte
}

func NewCoordProxy(parentIdentifier int, channel chan []byte) *CoordinatorProxy {
    return &CoordinatorProxy{parentIdentifier, channel}
}

func (self *CoordinatorProxy) RequestStatesInBox(turn int,
                                                 bottomLeft geo.Point,
                                                 topRight geo.Point) *GameStateResponse {
    
    self.conn <- GameStateRequestJson(self.parentIdentifier, turn, bottomLeft, topRight)
    
    timeout := time.NewTicker(5*1e9)
    select {
    case msg := <-self.conn:
        timeout.Stop()
        return GameStateResponseFromJson(msg)
    case <-timeout.C:
        timeout.Stop()
    }
    panic("RPC request timeout")
}

func (self *CoordinatorProxy) SendComplete() {
}
