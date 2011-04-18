package config

import geo "coord/geometry"

type AgentDefinition struct {
    Id uint
    X int
    Y int
}

func NewAgentDefinition(id uint, x, y int) *AgentDefinition {
    return &AgentDefinition{Id: id, X: x, Y: y}
}

type Config struct {
    Identifier int
    MaxTurns int
    Agents []*AgentDefinition
    MessageStyle string     // boolean|noise|none
    UseFood bool

    BottomLeft *geo.Point
    TopRight *geo.Point
}

func NewConfig(id int, maxTurns int, agents []*AgentDefinition, style string, food bool, bl, tr *geo.Point) *Config {
    return &Config{Identifier: id,
                   MaxTurns: maxTurns,
                   Agents: agents,
                   MessageStyle: style,
                   UseFood: food,
                   BottomLeft: bl,
                   TopRight: tr,
    }
}
