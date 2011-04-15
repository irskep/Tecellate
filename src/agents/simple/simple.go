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
    var x, y int
    if self.Id()%2 == 0 {
        x = 1
        y = 0
        comm.Log("broadcast success", comm.Broadcast(23, []byte("hello_world")))
    } else {
        comm.Log("broadcast success", comm.Broadcast(23, []byte("hello_world")))
        x = 1
        y = 1
    }
    if !comm.Move(x, y) {
        fmt.Println("Move Failed!")
    }
//     if !comm.Move(1, 0) {
//         fmt.Println("Move Failed!")
//     }
//     comm.Collect()
//     fmt.Println(comm.Look())
    comm.Log("listening", string(comm.Listen(23)))
    comm.Log("listening", string(comm.Listen(23)))
    comm.Log("listening", string(comm.Listen(23)))
//     fmt.Println(comm.PrevResult())
    comm.Log("my energy", comm.Energy())
    return
}

func (self *Simple) Id() uint {
    return self.id
}

func (self *Simple) Time() uint {
    return 0
}
