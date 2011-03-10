package coord

import geo "coord/geometry"

import (
    "time"
)

var timeout int64 = 5*1e9

type CoordinatorProxy struct {
    parentIdentifier int
    conn chan interface{}
}

func NewCoordProxy(parentIdentifier int, channel chan interface{}) *CoordinatorProxy {
    return &CoordinatorProxy{parentIdentifier, channel}
}

func (self *CoordinatorProxy) request(request interface{}) interface{} {
    self.conn <- request
    
    timeout := time.NewTicker(timeout)
    select {
    case response := <- self.conn:
        timeout.Stop()
        return response
    case <-timeout.C:
        timeout.Stop()
    }
    panic("RPC request timeout")
}

func (self *CoordinatorProxy) RequestStatesInBox(turn int,
                                                 bottomLeft geo.Point,
                                                 topRight geo.Point) GameStateResponse {
    request := GameStateRequest{self.parentIdentifier, turn, bottomLeft, topRight}
    return self.request(request).(GameStateResponse)
}

func (self *CoordinatorProxy) SendComplete() {
}
