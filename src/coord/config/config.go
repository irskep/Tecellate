package config

import geo "coord/geometry"

type AgentDefinition struct {
    Id uint32
    X int
    Y int
    Energy int
}

func NewAgentDefinition(id uint32, x, y, energy int) *AgentDefinition {
    return &AgentDefinition{Id: id, X: x, Y: y, Energy: energy}
}

type Config struct {
    Identifier int
    Address string
    MaxTurns int
    Agents []*AgentDefinition
    MessageStyle string     // boolean|noise|none
    UseFood bool

    BottomLeft *geo.Point
    TopRight *geo.Point
}

func NewConfig(id int, addr string, maxTurns int, agents []*AgentDefinition, style string, food bool, bl, tr *geo.Point) *Config {
    return &Config{Identifier: id,
                   Address: addr,
                   MaxTurns: maxTurns,
                   Agents: agents,
                   MessageStyle: style,
                   UseFood: food,
                   BottomLeft: bl,
                   TopRight: tr,
    }
}
