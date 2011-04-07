package coord

import geo "coord/geometry"
import cagent "coord/agent"

//     Turn() uint64
//     Position() *geo.Point
//     Energy() Energy
//     Alive() bool
//     Wait() uint16

type StateTransform struct {
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
                  wait uint16) *StateTransform {
    self := new(StateTransform)
    self.turn = turn
    self.pos = pos
    self.energy = energy
    self.alive = alive
    self.wait = wait
    return self
}

func (self *StateTransform) Turn() uint64 { return self.turn }
func (self *StateTransform) Position() *geo.Point { return self.pos }
func (self *StateTransform) Energy() cagent.Energy { return self.energy }
func (self *StateTransform) Alive() bool { return self.alive }
func (self *StateTransform) Wait() uint16 { return self.wait }
