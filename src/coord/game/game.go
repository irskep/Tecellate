package game

type GameState struct {
    Turn int
    Agents []*Agent
}

func (self *GameState) ApplyMoves(moves []*Move, agentStates []*AgentState) {
    
}
