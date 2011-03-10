/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agent/agent.go
*/

package agent

import "fmt"
import "agent/link"

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

type comm struct {
    conn link.Link
    done <-chan bool
}

func StartComm(conn link.Link) (*comm, chan<- bool) {
    self := new(comm)
    self.conn = conn
    done := make(chan bool)
    self.done = done
//     go func(self *comm) {
//         loop: for {
//             select {
//             case msg <- self.conn:
//                 //asdf
//             case <-self.done:
//                 break loop
//             }
//         }
//     }(self)
    return self, done
}

func (self *comm) Look() link.Vision {
    return nil
}

func (self *comm) Listen(freq uint8) link.Audio {
    return nil
}

func (self *comm) Broadcast(b link.Broadcast) bool {
    return false
}

func (self *comm) Inventory() link.Inventory {
    return nil
}

func (self *comm) Move(link.Move) bool {
    return false
}
func (self *comm) Collect() {
    return
}



func Run(agent Agent, conn link.Link) {
    complete := make(chan bool)
    go func(conn link.Link, done chan<- bool) {
        start := func() {
            cm, done := StartComm(conn)
            agent.Turn(cm)
            done <- true
        }

        loop: for {
            switch msg := <-conn; {
                case msg.Cmd == link.Commands["Start"]:
                    start()
                case msg.Cmd == link.Commands["Exit"]:
                    break loop
                default:
                    panic(
                        fmt.Sprintf(
                            "Command %s not valid for current state.",
                            msg.Cmd,
                        ),
                    )
            }
        }
        done <- true
    }(conn, complete)
    if ok := <-complete; ok {
        return
    }
    panic("we had an issue.")
}
