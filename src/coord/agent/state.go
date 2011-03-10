package agent

import geo "coord/geometry"

type AgentState struct {
    Turn uint64
    Live bool
    Position geo.Point
    Energy uint8
    Wait uint16  // the number of turns till the next movment

    NextMove Move
}

type Move struct {
    Position geo.Point
    Messages []Message
    Collect bool
}

type Message struct {
    Body string
    Frequency uint8
    Source geo.Point
}
