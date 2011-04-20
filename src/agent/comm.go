package agent

import (
    "fmt"
    "time"
    "logflow"
)
import (
    "agent/link"
    cagent "coord/agent"
    geo "coord/geometry"
)

type Comm interface {
    Log(...interface{})
    Logf(string, ...interface{})
    Look() link.Vision
    Listen(uint8) []byte
    Broadcast(uint8, []byte) bool
    Energy() cagent.Energy
    Move(x, y int) bool
    PrevResult() bool
    Collect() bool
}

type comm struct {
    snd link.SendLink
    rcv link.RecvLink
    log logflow.Logger
}

func StartComm(send link.SendLink, recv link.RecvLink, log logflow.Logger) *comm {
    self := new(comm)
    self.snd = send
    self.rcv = recv
    self.log = log
    return self
}

// func (self *comm) SwapChannels(newSnd...) {}

func (self *comm) Log(v ...interface{}) { self.log.Println(v...) }
func (self *comm) Logf(format string, v ...interface{}) { self.log.Printf(format, v...) }

func (self *comm) ack_start() {
    self.send(link.NewMessage(link.Commands["Ack"], link.Commands["Start"]))
}

func (self *comm) ack_migrate() {
    self.send(link.NewMessage(link.Commands["Ack"], link.Commands["Migrate"]))
}

func (self *comm) complete() bool {
    r, _ := self.acked_send(link.NewMessage(link.Commands["Complete"]))
    return r
}

func (self *comm) await_cmd_ack(cmd link.Command) (bool, link.Arguments) {
    msg := self.recv()

    proc := func(ack bool) (bool, link.Arguments) {
        acked := link.MakeCommand(msg.Args[0])
        if acked == cmd {
            if ack == true {
                return ack, msg.Args
            } else {
                return ack, nil
            }
        } else {
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

    if msg.Cmd == link.Commands["Ack"] && len(msg.Args) >= 1 {
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
//         self.log.Logln("proto", "recv :", msg)
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
//         self.log.Logln("proto", "recv :", msg)
        return &msg
    }
    panic("Did not recieve message.")
}

func (self *comm) send(msg *link.Message) {
    timeout := time.NewTicker(link.Timeout)
    select {
    case m := <-self.rcv:
//         self.log.Logln("proto", m)
        panic(fmt.Sprintf("unresolved message in pipe. \n msg = %v", m))
    case self.snd <- *msg:
//         self.log.Logln("proto", "sent :", msg)
    case <-timeout.C:
        timeout.Stop()
        panic("Agent believes the server to be unresponsive.")
    }
}

func (self *comm) acked_send(msg *link.Message) (bool, link.Arguments) {
    self.send(msg)
    return self.await_cmd_ack(msg.Cmd)
}

func (self *comm) Look() link.Vision {
    self.send(link.NewMessage(link.Commands["Look"]))
    self.recv()
    return nil
}

func (self *comm) Listen(freq uint8) []byte {
    m := link.NewMessage(link.Commands["Listen"], newListen(freq))
    if ok, args := self.acked_send(m); ok {
        if len(args) == 2 {
            return args[1]
        }
    }
    panic("didn't get an energy")
    self.recv()
    return nil
}

func (self *comm) Broadcast(freq uint8, msg []byte) bool {
    r, _ := self.acked_send(link.NewMessage(link.Commands["Broadcast"], newBroadcast(freq, msg)))
    return r
}

func (self *comm) Energy() cagent.Energy {
    if ok, args := self.acked_send(link.NewMessage(link.Commands["Energy"])); ok {
        if len(args) == 2 {
            return cagent.MakeEnergy(args[1])
        }
    }
    panic("didn't get an energy")
}

func (self *comm) Move(x,y int) bool {
    c, _ := self.acked_send(link.NewMessage(link.Commands["Move"], geo.NewPoint(x, y)))
    return c
}

func (self *comm) PrevResult() bool {
    self.send(link.NewMessage(link.Commands["PrevResult"]))
    self.recv()
    return false
}

func (self *comm) Collect() bool {
    r, _ := self.acked_send(link.NewMessage(link.Commands["Collect"]))
    return r
}
