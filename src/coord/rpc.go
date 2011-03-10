package coord

import geo "coord/geometry"

import (
    "coord/agent"
    "json"
)

type GameStateRequest struct {
    Turn int
    BottomLeft geo.Point
    TopRight geo.Point
}

func GameStateRequestJson(turn int, bottomLeft geo.Point, topRight geo.Point) []byte {
    requestBytes, _ :=  json.Marshal(GameStateRequest{turn, bottomLeft, topRight})
    return requestBytes
}

type GameStateResponse struct {
    Turn int
    AgentStates []agent.AgentState
}

func GameStateResponseJson(bytes []byte) *GameStateResponse {
    var response GameStateResponse
    _ = json.Unmarshal(bytes, &response)
    return &response
}
