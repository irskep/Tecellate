package agent

import "fmt"
import "time"
import "agent/link"

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

func (self *comm) ack_start() {
    self.conn <- *link.NewMessage(link.Commands["Ack"], link.Commands["Start"])
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

func (self *comm) Move(move link.Move) bool {
    self.conn <- *link.NewMessage(link.Commands["Move"], move)
    timeout := time.NewTicker(link.Timeout)
    select {
    case msg := <-self.conn:
        if msg.Cmd == link.Commands["Ack"] {
                return true
        } else if msg.Cmd == link.Commands["Nak"] {
                return false
        } else {
            s := fmt.Sprintf("Unexpected Message %s", msg)
            panic(s)
        }
    case <-timeout.C:
        timeout.Stop()
        panic("Server unresponsive")
    }
    panic("unreachable")
}

func (self *comm) Collect() {
    return
}

