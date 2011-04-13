package agent

import "fmt"
import geo "coord/geometry"

type Energy uint8

type AgentState struct {
    Id int
    Turn uint64
    Alive bool
    Position *geo.Point
    Energy Energy
    Wait uint16  // the number of turns till the next movment
    Move *Move
}

func NewAgentState(turn uint64, pos *geo.Point, energy Energy) *AgentState {
    self := &AgentState{
        Id:-1,
        Turn:turn,
        Alive:true,
        Position:pos,
        Wait:0,
        Energy:energy,
    }
    return self
}

func (self *AgentState) Transform(trans Transform) {
    self.Turn = trans.Turn()
    if trans.Position() != nil {
        self.Position = trans.Position()
    }
    self.Energy = trans.Energy()
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
    ok, _ := self.NewMessage(freq, msg)
    return ok
}

func (self *AgentState) PrevResult() bool {
    return false
}

func (self *AgentState) String() string {
    return fmt.Sprintf("<AgentState id:%v pos:%v>", self.Id, self.Position.String())
}
