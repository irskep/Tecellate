package coord

import geo "coord/geometry"

type CoordinatorProxy struct {
    conn chan []byte
}

func NewCoordProxyWithChannel(channel chan []byte) *CoordinatorProxy {
    return &CoordinatorProxy{channel}
}

func (self *CoordinatorProxy) RequestStatesInBox(turn int,
                                                 bottomLeft geo.Point,
                                                 topRight geo.Point) *GameStateResponse {
    self.conn <- GameStateRequestJson(turn, bottomLeft, topRight)
    return GameStateResponseFromJson(<- self.conn)
}

func (self *CoordinatorProxy) SendComplete() {
}
