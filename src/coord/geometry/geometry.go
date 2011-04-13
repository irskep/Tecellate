package geometry

import (
    "fmt"
    "math"
)

type Point struct {
    X int
    Y int
}

func NewPoint(x, y int) *Point {
    return &Point{X:x, Y:y}
}

func (self *Point) Add(other *Point) *Point {
    return &Point{
        X:self.X + other.X,
        Y:self.Y + other.Y,
    }
}

func (self *Point) Distance(other *Point) float64 {
    dx := float64(self.X - other.X)
    dy := float64(self.Y - other.Y)
    return math.Sqrt(dx*dx + dy*dy)
}
func (self *Point) DistanceSquare(other *Point) float64 {
    dx := float64(self.X - other.X)
    dy := float64(self.Y - other.Y)
    return dx*dx + dy*dy
}

func (self *Point) Complex() complex128 {
    return complex(float64(self.X), float64(self.Y))
}

func (self *Point) String() string {
    if self == nil {
        return "<nil>"
    }
    return fmt.Sprintf("<Point (%d, %d)>", self.X, self.Y)
}
