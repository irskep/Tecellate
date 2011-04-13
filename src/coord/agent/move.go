package agent

import "fmt"
import geo "coord/geometry"

type Move struct {
    Position *geo.Point
    Messages []*Message
    Collect bool
    setmv bool
}

func (self *AgentState) NewMove() *Move {
    self.Move = &Move{setmv:false}
    return self.Move
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
    return fmt.Sprintf("<Move %s %s>", self.Position.String(), self.Messages)
}
