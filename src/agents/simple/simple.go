/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agents/agent1.go
*/

package simple

import "agent"
import "fmt"

type Simple struct {
    id uint
}

func NewSimple(id uint) *Simple {
    return &Simple{id:id}
}

func (self *Simple) Turn(comm agent.Comm) {
    if !comm.Move(1, 0) {
        fmt.Println("Move Failed!")
    }
    if !comm.Move(1, 0) {
        fmt.Println("Move Failed!")
    }
    fmt.Println(comm.Look())
    fmt.Println(comm.Listen(23))
    fmt.Println(comm.Broadcast(23, []byte("hello")))
    return
}

func (self *Simple) Id() uint {
    return self.id
}
