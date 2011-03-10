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
}

func StartComm(conn link.Link) *comm {
    self := new(comm)
    self.conn = conn
    return self
}

func (self *comm) ack_start() {
    self.conn <- *link.NewMessage(link.Commands["Ack"], link.Commands["Start"])
}

func (self *comm) complete() bool {
    return self.acked_send(link.NewMessage(link.Commands["Complete"]))
}

func (self *comm) await_cmd_ack(cmd link.Command) bool {
    msg := self.recv()

    proc := func(ack bool) bool {
        switch acked := msg.Args[0].(type) {
        case link.Command:
            if acked == cmd {
                return ack
            }
        default:
            var s string
            if ack {
                s = fmt.Sprintf("Acked incorrect cmd (expected %s) %s", cmd, msg)
            } else {
                s = fmt.Sprintf("Ncked incorrect cmd (expected %s) %s", cmd, msg)
            }
            panic(s)
        }
        panic("unreachable")
    }

    if msg.Cmd == link.Commands["Ack"] && len(msg.Args) == 1 {
        return proc(true)
    } else if msg.Cmd == link.Commands["Nak"] && len(msg.Args) == 1 {
        return proc(false)
    }
    s := fmt.Sprintf("Unexpected Message %s", msg)
    panic(s)
}

func (self *comm) recv() *link.Message {
    timeout := time.NewTicker(link.Timeout)
    select {
    case msg := <-self.conn:
        return &msg
    case <-timeout.C:
        timeout.Stop()
        panic("Agent believes the server to be unresponsive.")
    }
    panic("Did not recieve message.")
}

func (self *comm) send(msg *link.Message) {
    timeout := time.NewTicker(link.Timeout)
    select {
    case self.conn <- *msg:
    case <-timeout.C:
        timeout.Stop()
        panic("Agent believes the server to be unresponsive.")
    }
}

func (self *comm) acked_send(msg *link.Message) bool {
    self.send(msg)
    return self.await_cmd_ack(msg.Cmd)
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
    return self.acked_send(link.NewMessage(link.Commands["Move"], move))
}

func (self *comm) Collect() {
    return
}

