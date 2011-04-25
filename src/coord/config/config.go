package config

import (
    geo "coord/geometry"
    "logflow"
)

type LogConfig []string
type LogConfigList []LogConfig

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
    
    Logs LogConfigList
    Peers map[string]int

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
                   Logs: make(LogConfigList, 0),
                   Peers: make(map[string]int),
                   BottomLeft: bl,
                   TopRight: tr,
    }
}

// Logs

func (self LogConfigList) Apply() {
    for _, l := range(self) {
        l.Apply()
    }
}

func (self LogConfig) Apply() {
    switch self[0] {
    case "stdout":
        logflow.StdoutSink(self[1])
    case "file":
        logflow.FileSink(self[1], true, self[2:]...)
    }
}
