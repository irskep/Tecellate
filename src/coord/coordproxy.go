package coord

import "coord.game"

type CoordinatorProxy struct {
    conn chan []byte
}

type GameStateResponse struct {
    Turn int
    AgentStates []AgentState
}

func (self *CoordinatorProxy) RequestStatesInBox(bottomLeft Point, topRight Point, turn int) *GameStateResponse {
    return &GameStateResponse{turn, nil};
}

func (self *CoordinatorProxy) SendComplete() {
}
