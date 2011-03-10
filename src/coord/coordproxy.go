package coord

import geo "coord/geometry"

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
    return GameStateResponseFromJson(<- self.conn)
}

func (self *CoordinatorProxy) SendComplete() {
}
