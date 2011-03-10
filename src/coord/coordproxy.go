package coord

import geo "coord/geometry"

import (
    "coord/agent"
    "json"
)

type CoordinatorProxy struct {
    conn chan []byte
}

type GameStateRequest struct {
    Turn int
    BottomLeft geo.Point
    TopRight geo.Point
}

type GameStateResponse struct {
    Turn int
    AgentStates []agent.AgentState
}

func NewCoordProxyWithChannel(channel chan []byte) *CoordinatorProxy {
    return &CoordinatorProxy{channel}
}

func (self *CoordinatorProxy) RequestStatesInBox(turn int,
                                                 bottomLeft geo.Point,
                                                 topRight geo.Point) *GameStateResponse {
    requestBytes, _ :=  json.Marshal(GameStateRequest{turn, bottomLeft, topRight})
    self.conn <- requestBytes
    var response GameStateResponse
    _ = json.Unmarshal(<- self.conn, &response)
    // TODO: Check err
    return &response
}

func (self *CoordinatorProxy) SendComplete() {
}
