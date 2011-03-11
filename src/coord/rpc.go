package coord

import geo "coord/geometry"

import (
    "coord/agent"
)

type GameStateRequest struct {
    SenderIdentifier int
    Turn int
    BottomLeft geo.Point
    TopRight geo.Point
}

type GameStateResponse struct {
    Turn int
    AgentStates []agent.AgentState
}

func (self GameStateResponse) CopyToHeap() *GameStateResponse {
    return &GameStateResponse{self.Turn, self.AgentStates}
}
