/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agent/agent.go
*/

package agent

import "agent/link"
import geo "coord/geometry"

type Agent interface {
    Turn(Comm)
}

type Comm interface {
    Look() link.Vision
    Listen(uint8) link.Audio
    Broadcast(link.Broadcast) bool
    Inventory() link.Inventory
    Move(link.Move) bool
    Collect()
}

type Move struct {
    pos geo.Point
}

func NewMove(x, y int) *Move {
    self := new(Move)
    self.pos = *geo.NewPoint(x,y)
    return self
}
