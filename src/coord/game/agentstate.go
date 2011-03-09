package game

type AgentState struct {
    Turn int
    Live bool
    Position Point
}

type AgentProxy struct {
    State AgentState
}

type Agent interface {
    Turn()
}

func (self *AgentProxy) Turn() {
    return
}
