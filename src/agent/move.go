package agent

import "fmt"
import geo "coord/geometry"

type Move struct {
    pos geo.Point
}

func newMove(x, y int) *Move {
    self := new(Move)
    self.pos = *geo.NewPoint(x,y)
    return self
}

func (self *Move) Move() geo.Point {
    return self.pos
}

func (self *Move) String() string {
    return fmt.Sprintf("%d, %d", self.pos.X, self.pos.Y)
}
