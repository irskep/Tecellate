package config

import geo "coord/geometry"
import "coord/agent"

type Config struct {
    Identifier int
    MaxTurns int
    Agents []agent.Agent
    MessageStyle string     // boolean|noise|none
    UseFood bool
    RandomlyDelayProcessing bool

    BottomLeft *geo.Point
    TopRight *geo.Point
}

func NewConfig(id int, maxTurns int, agents []agent.Agent, style string, food bool, delay bool, bl, tr *geo.Point) *Config {
    return &Config{Identifier: id,
                   MaxTurns: maxTurns,
                   Agents: agents,
                   MessageStyle: style,
                   UseFood: food,
                   RandomlyDelayProcessing: delay,
                   BottomLeft: bl,
                   TopRight: tr,
    }
}
