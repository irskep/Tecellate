package agent

import geo "coord/geometry"

// State

type AgentState struct {
    Turn int
    Live bool
    Position geo.Point
    Energy int
    TurnsToNextAllowedMove uint

    NextMove Move
}

type Move struct {
    Position geo.Point
    Messages []Message
    Collect bool
}

type Message struct {
    Body string
    Frequency int
    Source geo.Point
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
