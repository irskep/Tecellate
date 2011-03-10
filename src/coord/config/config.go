package config

import geo "coord/geometry"

type Config struct {
    AgentStarts []AgentStart
}

type AgentStart struct {
    Position geo.Point
    Kind string
}
