package game

import "coord/agent"
import "coord/config"

type GameState struct {
    Turn uint64
    Agents []agent.Agent
    Terrain *Map
    Energy *Map
    conf *config.Config
}

func NewGameState() *GameState {
    return &GameState{0, make([]agent.Agent, 0), nil, nil, nil}
}

func (self *GameState) Advance() {
    self.Turn += 1
}

func (self *GameState) Copy() *GameState {
    return &GameState{
            self.Turn,
            self.Agents,
            self.Terrain.Copy(),
            self.Energy.Copy(),
            self.conf,
    }
}

func (self *GameState) Configure(conf *config.Config) {
    self.conf = conf
    self.Agents = conf.Agents
}

func (self *GameState) ApplyMoves(moves []*agent.Move, agentStates []*agent.AgentState) {

}

type Map struct {
    Values [][]int
    Width uint
    Height uint
}

func NewMap(w uint, h uint) *Map {
    return &Map{make([][]int, w, h), w, h}
}

func (self *Map) Copy() *Map {
    newMap := NewMap(self.Width, self.Height)
    for i := uint(0); i < self.Width; i++ {
        for j := uint(0); j < self.Height; j++ {
            newMap.Values[i][j] = self.Values[i][j]
        }
    }
    return newMap
}
