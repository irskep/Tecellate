package config

import geo "coord/geometry"

type Config struct {
    AgentStarts []AgentStart
    MessageStyle string     // boolean|noise|none
    UseFood bool
}

type AgentStart struct {
    Position geo.Point
    Kind string
}
