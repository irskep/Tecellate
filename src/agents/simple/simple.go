/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agents/agent1.go
*/

package simple

import "agent"

type Simple struct {
    /// pass
}

func NewSimple() *Simple {
    return &Simple{}
}

func (self *Simple) Turn(comm agent.Comm) {
    comm.Move(agent.NewMove(1, 0))
    return
}
