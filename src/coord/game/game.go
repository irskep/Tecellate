package game

import cagent "coord/agent"

import "coord/config"

import (
    "fmt"
)

type GameState struct {
    Turn uint64
    Agents []cagent.Agent
    Terrain *Map
    Energy *Map
    conf *config.Config
    statesToServe []cagent.AgentState
}

func NewGameState() *GameState {
    return &GameState{0, make([]cagent.Agent, 0), nil, nil, nil, nil}
}

func (self *GameState) Advance(transforms []cagent.Transform) {
    self.Turn += 1
    self.statesToServe = nil
    for i, agent := range(self.Agents) {
        agent.Apply(transforms[i])
    }
}

func (self *GameState) Configure(conf *config.Config) {
    self.conf = conf
    self.Agents = conf.Agents
}

func (self *GameState) AgentStates() []cagent.AgentState {
    if self.statesToServe == nil {
        self.statesToServe = make([]cagent.AgentState, len(self.Agents))
        for _, agent := range(self.Agents) {
            self.statesToServe = append(self.statesToServe, *agent.State())
        }
    }
    return self.statesToServe
}

func (self *GameState) MakeRPCResponse() GameStateResponse {
    return GameStateResponse{self.Turn, self.AgentStates()}
}

// RPC response

type GameStateResponse struct {
    Turn uint64
    AgentStates []cagent.AgentState
}

func (self GameStateResponse) CopyToHeap() *GameStateResponse {
    return &GameStateResponse{self.Turn, self.AgentStates}
}

func (self GameStateResponse) String() string {
    
    return fmt.Sprintf("Turn %d: %v", self.Turn, self.AgentStates)
}
