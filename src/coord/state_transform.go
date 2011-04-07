package coord

import geo "coord/geometry"
import cagent "coord/agent"

//     Turn() uint64
//     Position() *geo.Point
//     Energy() Energy
//     Alive() bool
//     Wait() uint16

type Transform struct {
    turn uint64
    pos *geo.Point
    energy cagent.Energy
    alive bool
    wait uint16
}

func newTransform(turn uint64,
                  pos *geo.Point,
                  energy cagent.Energy,
                  alive bool,
                  wait uint16) *Transform {
    self := new(Transform)
    self.turn = turn
    self.pos = pos
    self.energy = energy
    self.alive = alive
    self.wait = wait
    return self
}

func (self *Transform) Turn() uint64 { return self.turn }
func (self *Transform) Position() *geo.Point { return self.pos }
func (self *Transform) Energy() cagent.Energy { return self.energy }
func (self *Transform) Alive() bool { return self.alive }
func (self *Transform) Wait() uint16 { return self.wait }
