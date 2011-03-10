package geometry

import (
    "math"
)

type Point struct {
    X int
    Y int
}

func (self Point) Distance(other Point) float64 {
    dx := float64(self.X - other.X)
    dy := float64(self.Y - other.Y)
    return math.Sqrt(dx*dx + dy*dy)
}
