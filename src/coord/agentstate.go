package coord

type AgentState struct {
    Turn int
    Live bool
    
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
