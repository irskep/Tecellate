package config

import geo "coord/geometry"

type Config struct {
    Identifier int
    AgentStarts []AgentStart
    MessageStyle string     // boolean|noise|none
    UseFood bool
    RandomlyDelayProcessing bool
}

type AgentStart struct {
    Position geo.Point
    Kind string
}
