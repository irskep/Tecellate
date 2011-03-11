package agent

import geo "coord/geometry"

type Energy uint8

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

type Transform interface {
    Turn() uint64
    Position() *geo.Point
    Energy() Energy
    Alive() bool
    Wait() uint16
}

type Message struct {
    Msg []byte
    Frequency uint8
    Source *geo.Point
}

type Inventory struct {
    Energy Energy
}

func NewAgentState(turn uint64, pos *geo.Point, energy Energy) *AgentState {
    self := &AgentState{
        Turn:turn,
        Live:true,
        Position:pos,
        Wait:0,
        Inventory:NewInventory(energy),
    }
    return self
}

func NewInventory(energy Energy) *Inventory {
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

func (self *AgentState) transform(trans Transform) {
    self.Turn = trans.Turn()
    self.Position = trans.Position()
    self.Inventory.Energy = trans.Energy()
    self.Live = trans.Alive()
    self.Wait = trans.Wait()
}

func (self *AgentState) Mv(pos *geo.Point) bool {
    return self.Move.mv(pos)
}

func (self *Move) mv(pos *geo.Point) bool {
    if !self.setmv {
        self.Position = pos
        self.setmv = true
        return true
    }
    return false
}
