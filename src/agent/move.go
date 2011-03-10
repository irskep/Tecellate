package agent


import geo "coord/geometry"

type Move struct {
    pos geo.Point
}

func NewMove(x, y int) *Move {
    self := new(Move)
    self.pos = *geo.NewPoint(x,y)
    return self
}

func (self *Move) Move() geo.Point {
    return self.pos
}
