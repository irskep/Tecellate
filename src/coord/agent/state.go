package agent

import "fmt"
import geo "coord/geometry"

type Energy uint8

type AgentState struct {
    Id int
    Turn uint64
    Alive bool
    Position *geo.Point
    Inventory *Inventory
    Wait uint16  // the number of turns till the next movment
    Move *Move
}

type Inventory struct {
    Energy Energy
}

func NewAgentState(turn uint64, pos *geo.Point, energy Energy) *AgentState {
    self := &AgentState{
        Id:-1,
        Turn:turn,
        Alive:true,
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

func (self *AgentState) transform(trans Transform) {
    self.Turn = trans.Turn()
    if trans.Position() != nil {
        self.Position = trans.Position()
    }
    self.Inventory.Energy = trans.Energy()
    self.Alive = trans.Alive()
    self.Wait = trans.Wait()
}

func (self *AgentState) Mv(pos *geo.Point) bool {
    return self.Move.mv(pos)
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

func (self *AgentState) String() string {
    return fmt.Sprintf("<AgentState id:%v pos:%v>", self.Id, self.Position.String())
}
