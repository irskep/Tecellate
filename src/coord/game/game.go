package game

type GameState struct {
    Turn int
    Agents []*Agent
}

func NewGameState() *GameState {
    return &GameState{0, make([]*Agent, 0)}
}

func (self *GameState) Configure(config Config) {
    
}

func (self *GameState) ApplyMoves(moves []*Move, agentStates []*AgentState) {
    
}
