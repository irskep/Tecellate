package game

import (
    . "byteslice"
    geo "coord/geometry"
    cagent "coord/agent"
    "coord/config"
    "fmt"
    "logflow"
)

type GameState struct {
    Turn uint64
    Agents []cagent.Agent
    Terrain *Map
    Energy *Map
    conf *config.Config
    statesToServe []cagent.AgentState
    messages *Messages
    myMessages *Messages
    upForAdoption []cagent.Agent
}

type GameStateRequest struct {
    SenderIdentifier int
    SenderAddress string
    Turn int
    BottomLeft geo.Point
    TopRight geo.Point
}

type GameStateResponse struct {
    Identifier int
    Turn uint64
    AgentStates []cagent.AgentState
    AgentsToAdopt []cagent.AgentState
    Messages map[uint8][]cagent.Message
}

func (self *GameStateResponse) CopyToHeap() *GameStateResponse {
    return &GameStateResponse{self.Identifier, self.Turn, self.AgentStates, self.AgentsToAdopt, self.Messages}
}

func (self *GameStateResponse) String() string {
    return fmt.Sprintf("Turn %d: %v (%d messages)", self.Turn, self.AgentStates, len(self.Messages))
}

func NewGameState() *GameState {
    return &GameState{
        Turn:0,
        Agents:make([]cagent.Agent, 0),
        messages:NewMessages(nil),
        myMessages:NewMessages(nil),
        upForAdoption:make([]cagent.Agent, 0),
    }
}

func (self *GameState) Advance(transforms []cagent.Transform, messages *Messages, myMessages *Messages) {
    self.Turn += 1
    
    for i, agent := range(self.Agents) {
        agent.Apply(transforms[i])
    }
    
    self.statesToServe = make([]cagent.AgentState, len(self.Agents))
    for i, agent := range(self.Agents) {
        self.statesToServe[i] = *agent.State()
    }
    
    bl := self.conf.BottomLeft
    tr := self.conf.TopRight
    allNewAgents := make([]cagent.Agent, 0)
    self.upForAdoption = make([]cagent.Agent, 0)
    for _, agent := range(self.Agents) {
        p := agent.State().Position
        if bl.X <= p.X && bl.Y <= p.Y && p.X < tr.X && p.Y < tr.Y {
            allNewAgents = append(allNewAgents, agent)
        } else {
            self.upForAdoption = append(self.upForAdoption, agent)
        }
    }
    
    self.Agents = allNewAgents
    self.messages = messages
    self.myMessages = myMessages
}

func (self *GameState) Configure(conf *config.Config) {
    self.conf = conf
}

func (self *GameState) AgentStates() []cagent.AgentState {
    return self.statesToServe
}

func (self *GameState) AgentStatesToExport(req GameStateRequest) []cagent.AgentState {
    bl := req.BottomLeft
    tr := req.TopRight
    upForAdoption := make([]cagent.AgentState, 0)
    for _, agent := range(self.upForAdoption) {
        p := agent.State().Position
        if bl.X <= p.X && bl.Y <= p.Y && p.X < tr.X && p.Y < tr.Y {
            upForAdoption = append(upForAdoption, *agent.State())
            agent.MigrateTo(req.SenderAddress)
        }
    }
    if len(upForAdoption) > 0 {
        logflow.Print("gamestate/info", "Transferring from ", self.conf.Identifier, " to ", req.SenderIdentifier, ":")
        logflow.Print("gamestate/info", upForAdoption)
    }
    return upForAdoption
}

func (self *GameState) MakeRPCResponse(req GameStateRequest) GameStateResponse {
    return GameStateResponse{self.conf.Identifier,
                             self.Turn, self.AgentStates(), 
                             self.AgentStatesToExport(req),
                             self.myMessages.Msgs}
}

func (self *GameState) Listen(loc geo.Point, freq uint8) ByteSlice {
    return self.messages.Hear(loc, freq)
}
