package agent

import geo "coord/geometry"

import (
    "fmt"
)

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

func (self *AgentState) String() string {
    return fmt.Sprintf("AgentState(turn=%u, live=%t, position=%v, inventory=%v, wait=%u, move=%v)",
                       self.Turn, self.Live, self.Position.String(), self.Inventory, self.Wait, self.Move.String())
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

func (self *Move) String() string {
    if self == nil {
        return "<nil>"
    }
    return fmt.Sprintf("Move(position=%v, %d messages, collect=%t", 
                       self.Position.String(),
                       len(self.Messages),
                       self.Collect)
}

func (self *AgentState) Collect() bool {
    return false
}

func (self *AgentState) Listen(freq uint8) []byte {
    return nil
}

func (self *AgentState) Broadcast(freq uint8, msg []byte) bool {
    return false
}

func (self *AgentState) PrevResult() bool {
    return false
}

func (self *AgentState) GetInventory() *Inventory {
    return self.Inventory
}
