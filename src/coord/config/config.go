package config

import geo "coord/geometry"
import "coord/agent"

type Config struct {
    Identifier int
    Agents []agent.Agent
    MessageStyle string     // boolean|noise|none
    UseFood bool
    RandomlyDelayProcessing bool

    BottomLeft *geo.Point
    TopRight *geo.Point
}

func NewConfig(id int, agents []agent.Agent, style string, food bool, delay bool, bl, tr *geo.Point) *Config {
    return &Config{Identifier: id,
                   Agents: agents,
                   MessageStyle: style,
                   UseFood: food,
                   RandomlyDelayProcessing: delay,
                   BottomLeft: bl,
                   TopRight: tr,
    }
}

func (self *Config) Duplicate(identifier int, bottomLeft, topRight *geo.Point) *Config {
    return &Config{Identifier: identifier,
                   Agents: self.Agents,
                   MessageStyle: self.MessageStyle,
                   UseFood: self.UseFood,
                   RandomlyDelayProcessing: self.RandomlyDelayProcessing,
                   BottomLeft: bottomLeft,
                   TopRight: topRight,
    }
}
