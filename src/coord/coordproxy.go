package coord

import geo "coord/geometry"
import "coord/agent"

type CoordinatorProxy struct {
    conn chan []byte
}

type GameStateResponse struct {
    Turn int
    AgentStates []agent.AgentState
}

func NewCoordProxyWithChannel(channel chan []byte) *CoordinatorProxy {
    return &CoordinatorProxy{channel}
}

func (self *CoordinatorProxy) RequestStatesInBox(bottomLeft geo.Point,
                                                 topRight geo.Point,
                                                 turn int) *GameStateResponse {
    return &GameStateResponse{turn, nil};
}

func (self *CoordinatorProxy) SendComplete() {
}
