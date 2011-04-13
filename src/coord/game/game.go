package game

import (
    "fmt"
)

import (
    geo "coord/geometry"
    cagent "coord/agent"
    "coord/config"
)


type GameState struct {
    Turn uint64
    Agents []cagent.Agent
    Terrain *Map
    Energy *Map
    conf *config.Config
    statesToServe []cagent.AgentState
    messages *Messages
}

func NewGameState() *GameState {
    return &GameState{
        Turn:0,
        Agents:make([]cagent.Agent, 0),
        messages:NewMessages(nil),
    }
}

func (self *GameState) Advance(transforms []cagent.Transform, messages *Messages) {
    self.Turn += 1
    self.statesToServe = nil
    for i, agent := range(self.Agents) {
        agent.Apply(transforms[i])
    }
    self.messages = messages
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
    return GameStateResponse{self.Turn, self.AgentStates(), self.messages.Msgs}
}

func (self *GameState) Listen(loc *geo.Point, freq uint8) []byte {
    return self.messages.Hear(loc, freq)
}

// RPC response

type GameStateResponse struct {
    Turn uint64
    AgentStates []cagent.AgentState
    Messages map[uint8][]cagent.Message
}

func (self GameStateResponse) CopyToHeap() *GameStateResponse {
    return &GameStateResponse{self.Turn, self.AgentStates, self.Messages}
}

func (self GameStateResponse) String() string {
    return fmt.Sprintf("Turn %d: %v (%d messages)", self.Turn, self.AgentStates, len(self.Messages))
}
