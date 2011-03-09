package game

// State

type AgentState struct {
    Turn int
    Live bool
    Position Point
    Energy int
    TurnsToNextAllowedMove uint
    
    NextMove Move
}

type Move struct {
    Position Point
    Messages []Message
    Collect bool
}

type Message struct {
    Body string
    Frequency int
    Source Point
}

// Actions

type AgentProxy struct {
    State AgentState
    conn chan []byte
}

type Agent interface {
    Turn(completionBlock chan bool)
}

func (self *AgentProxy) Turn(completionBlock chan bool) {
    if (completionBlock != nil) {
        completionBlock <- true
    }
}
