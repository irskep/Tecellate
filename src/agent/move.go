package agent

import "fmt"
import geo "coord/geometry"

type move struct {
    pos geo.Point
}

func newMove(x, y int) *move {
    self := new(move)
    self.pos = *geo.NewPoint(x,y)
    return self
}

func (self *move) Move() geo.Point {
    return self.pos
}

func (self move) String() string {
    return fmt.Sprintf("%d, %d", self.pos.X, self.pos.Y)
}
