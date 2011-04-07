package agent

import (
    "fmt"
    "time"
    "log"
)
import (
    "agent/link"
)

type Comm interface {
    Look() link.Vision
    Listen(uint8) link.Audio
    Broadcast(uint8, []byte) bool
    GetInventory() link.Inventory
    Move(x, y int) bool
    PrevResult() bool
    Collect() bool
}

type comm struct {
    snd link.SendLink
    rcv link.RecvLink
    log *log.Logger
}

func StartComm(send link.SendLink, recv link.RecvLink, log *log.Logger) *comm {
    self := new(comm)
    self.snd = send
    self.rcv = recv
    self.log = log
    return self
}

func (self *comm) ack_start() {
    self.send(link.NewMessage(link.Commands["Ack"], link.Commands["Start"]))
}

func (self *comm) id(id uint) {
    self.send(link.NewMessage(link.Commands["Ack"], link.Commands["Id"], id))
}

func (self *comm) complete() bool {
//     fmt.Println("started complete")
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
    case msg := <-self.rcv:
        self.log.Println("recv :", msg)
        return &msg
    case <-timeout.C:
        timeout.Stop()
        panic("Agent believes the server to be unresponsive.")
    }
    panic("Did not recieve message.")
}

func (self *comm) recv_forever() *link.Message {
    select {
    case msg := <-self.rcv:
        self.log.Println("recv :", msg)
        return &msg
    }
    panic("Did not recieve message.")
}

func (self *comm) send(msg *link.Message) {
    timeout := time.NewTicker(link.Timeout)
    select {
    case m := <-self.rcv:
        self.log.Println(m)
        panic("unresolved message in pipe.")
    case self.snd <- *msg:
        self.log.Println("sent :", msg)
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
    self.send(link.NewMessage(link.Commands["Look"]))
    self.recv()
    return nil
}

func (self *comm) Listen(freq uint8) link.Audio {
    self.send(link.NewMessage(link.Commands["Listen"], newListen(freq)))
    self.recv()
    return nil
}

func (self *comm) Broadcast(freq uint8, msg []byte) bool {
    return self.acked_send(link.NewMessage(link.Commands["Broadcast"], newBroadcast(freq, msg)))
}

func (self *comm) GetInventory() link.Inventory {
    self.send(link.NewMessage(link.Commands["Inventory"]))
    self.recv()
    return nil
}

func (self *comm) Move(x,y int) bool {
    c := self.acked_send(link.NewMessage(link.Commands["Move"], newMove(x, y)))
//     fmt.Println("completed move")
    return c
}

func (self *comm) PrevResult() bool {
    self.send(link.NewMessage(link.Commands["PrevResult"]))
    self.recv()
    return false
}

func (self *comm) Collect() bool {
    return self.acked_send(link.NewMessage(link.Commands["Collect"]))
}
