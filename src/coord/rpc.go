// TODO: CHECK ANY ERRORS WHATSOEVER

package coord

import geo "coord/geometry"

import (
    "coord/agent"
    "json"
)

/* Request */

type GameStateRequest struct {
    SenderIdentifier int
    Turn int
    BottomLeft geo.Point
    TopRight geo.Point
}

func GameStateRequestJson(identifier int, turn int, bottomLeft geo.Point, topRight geo.Point) []byte {
    requestBytes, _ :=  json.Marshal(GameStateRequest{identifier, turn, bottomLeft, topRight})
    return requestBytes
}

func GameStateRequestFromJson(bytes []byte) *GameStateRequest {
    var request GameStateRequest
    _ = json.Unmarshal(bytes, &request)
    return &request
}

/* Response */

type GameStateResponse struct {
    Turn int
    AgentStates []agent.AgentState
}

func GameStateResponseJson(turn int, agentStates []agent.AgentState) []byte {
    responseBytes, _ :=  json.Marshal(GameStateResponse{turn, agentStates})
    return responseBytes
}

func GameStateResponseFromJson(bytes []byte) *GameStateResponse {
    var response GameStateResponse
    _ = json.Unmarshal(bytes, &response)
    return &response
}
