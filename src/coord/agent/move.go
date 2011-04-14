package agent

import "fmt"
import geo "coord/geometry"

type Move struct {
    Valid bool
    Position geo.Point
    Messages []Message
    Collect bool
    Setmv bool
}

func (self *AgentState) NewMove() Move {
    self.Move = Move{Valid:true, Setmv:false}
    return self.Move
}

func (self *Move) mv(pos geo.Point) bool {
    if !self.Setmv {
        self.Position = pos
        self.Setmv = true
        return true
    }
    return false
}

func (self *Move) String() string {
    if self == nil {
        return "<nil>"
    }
    return fmt.Sprintf("<Move %s %s>", self.Position.String(), self.Messages)
}
