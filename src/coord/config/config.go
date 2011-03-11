package config

import geo "coord/geometry"

type Config struct {
    Identifier int
    AgentStarts []AgentStart
    MessageStyle string     // boolean|noise|none
    UseFood bool
    RandomlyDelayProcessing bool
    
    BottomLeft geo.Point
    TopRight geo.Point
}

func BasicTestConfig() *Config {
    return &Config{Identifier: 0, 
                   AgentStarts: nil, 
                   MessageStyle: "none", 
                   UseFood: false, 
                   RandomlyDelayProcessing: true}
}

func (self *Config) Duplicate(identifier int, bottomLeft geo.Point, topRight geo.Point) *Config {
    return &Config{Identifier: identifier, 
                   AgentStarts: self.AgentStarts, 
                   MessageStyle: self.MessageStyle, 
                   UseFood: self.UseFood, 
                   RandomlyDelayProcessing: self.RandomlyDelayProcessing,
                   BottomLeft: bottomLeft,
                   TopRight: topRight}
}

type AgentStart struct {
    Position geo.Point
    Kind string
}
