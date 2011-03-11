package agent

import geo "coord/geometry"

type AgentState struct {
    Turn uint64
    Live bool
    Position *geo.Point
    Inventory *Inventory
    Wait uint16  // the number of turns till the next movment
    Move *Move
}

type Move struct {
    Position *geo.Point
    Messages []*Message
    Collect bool
    setmv bool
}

type Message struct {
    Msg []byte
    Frequency uint8
    Source *geo.Point
}

type Inventory struct {
    Energy uint8
}

func NewAgentState(turn uint64, pos *geo.Point, energy uint8) *AgentState {
    return &AgentState{
        Turn:turn,
        Live:true,
        Position:pos,
        Wait:0,
        Inventory:NewInventory(energy),
    }
}

func NewInventory(energy uint8) *Inventory {
    return &Inventory{
        Energy:energy,
    }
}

func (self *AgentState) NewMove() *Move {
    self.Move = &Move{setmv:false}
    return self.Move
}

func (self *AgentState) NewMessage(freq uint8, msg []byte) *Message {
    m := &Message{Frequency:freq, Msg:msg, Source:self.Position}
    self.Move.Messages = append(self.Move.Messages, m)
    return m
}

func (self *Move) Move(pos *geo.Point) bool {
    if self.setmv {
        self.Position = pos
        self.setmv = true
        return true
    }
    return false
}
