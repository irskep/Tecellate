package coord

import "coord/game"

type CoordinatorProxy struct {
    conn chan []byte
}

type GameStateResponse struct {
    Turn int
    AgentStates []game.AgentState
}

func NewFromLocal(other *Coordinator) *CoordinatorProxy {
    return nil;
}

func (self *CoordinatorProxy) RequestStatesInBox(bottomLeft game.Point, 
                                                 topRight game.Point, 
                                                 turn int) *GameStateResponse {
    return &GameStateResponse{turn, nil};
}

func (self *CoordinatorProxy) SendComplete() {
}
