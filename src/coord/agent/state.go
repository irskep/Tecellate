package agent

import "fmt"
import geo "coord/geometry"
import . "byteslice"

type Energy uint32
type Turn uint64

type AgentState struct {
    Id uint32
    Turn Turn
    Alive bool
    Position geo.Point
    Energy Energy
    Wait uint16  // the number of turns till the next movment
    Move Move
}

func MakeEnergy(bytes ByteSlice) Energy {
    return Energy(bytes.Int32())
}

func (self Energy) Bytes() ByteSlice {
    return ByteSlice32(uint32(self))
}



func MakeTurn(bytes ByteSlice) Turn {
    return Turn(bytes.Int64())
}

func (self Turn) Bytes() ByteSlice {
    return ByteSlice64(uint64(self))
}

func NewAgentState(id uint32, turn uint64, pos geo.Point, energy Energy) *AgentState {
    self := &AgentState{
        Id:id,
        Turn:Turn(turn),
        Alive:true,
        Position:pos,
        Wait:0,
        Energy:energy,
    }
    return self
}

func (self *AgentState) Transform(trans Transform) {
    self.Turn = Turn(trans.Turn())
    self.Position = trans.Position()
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

func (self AgentState) String() string {
    return fmt.Sprintf("<AgentState id:%v pos:%v>", self.Id, self.Position.String())
}
