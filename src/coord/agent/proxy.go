package agent

import (
    "fmt"
    "time"
    "log"
    "os"
)
import (
    "agent/link"
//     geo "coord/geometry"
)

type AgentProxy struct {
    state *AgentState
    snd link.SendLink
    rcv link.RecvLink
    log *log.Logger
}

func NewAgentProxy(send link.SendLink, recv link.RecvLink) *AgentProxy {
    self := new(AgentProxy)
//     self.state = NewAgentState(0, geo.NewPoint(0, 0), 0)
    self.snd = send
    self.rcv = recv
    self.log = log.New(os.Stdout, "AgentProxy : ", 0)
    return self
}

func (self *AgentProxy) SetState(state *AgentState) {
    self.state = state
}

func (self *AgentProxy) State() *AgentState {
    return self.state
}

func (self *AgentProxy) Apply(trans Transform) {
    self.state.transform(trans)
}

func (self *AgentProxy) Turn() bool {
    type handler (func(*link.Message) bool)

    check_args := func(count int, args link.Arguments) bool {
        if len(args) == count {
            return true
        }
        self.log.Println("Error : Wrong number of arguments recieved")
        return false
    }

    argnum := func(count int, f handler) handler {
        return func(msg *link.Message) bool {
            if check_args(count, msg.Args) {
                return f(msg)
            }
            return false
        }
    }

    var handlers = map[link.Command]handler {
        link.Commands["Complete"]:
            argnum(0, func(msg *link.Message) bool {
                self.ack_cmd(msg.Cmd)
                return true
            }),
        link.Commands["Move"]:
            argnum(1, func(msg *link.Message) bool {
                mv := msg.Args[0].(link.Move).Move()
                if self.state.Mv(&mv) {
                    self.ack_cmd(msg.Cmd)
                } else {
                    self.nak_cmd(msg.Cmd)
                }
                return false
            }),
        link.Commands["Look"]:
            argnum(0, func(msg *link.Message) bool {
                self.send(link.NewMessage(link.Commands["Ack"], msg.Cmd, nil))
                return false
            }),
        link.Commands["Listen"]:
            argnum(1, func(msg *link.Message) bool {
                freq := msg.Args[0].(link.Listen).Listen()
                self.send(link.NewMessage(link.Commands["Ack"], msg.Cmd, freq))
                return false
            }),
        link.Commands["Broadcast"]:
            argnum(1, func(msg *link.Message) bool {
                freq, pkt := msg.Args[0].(link.Broadcast).Message()
                if self.state.Broadcast(freq, pkt) {
                    self.ack_cmd(msg.Cmd)
                } else {
                    self.nak_cmd(msg.Cmd)
                }
                return false
            }),
        link.Commands["Collect"]:
            argnum(0, func(msg *link.Message) bool {
                if self.state.Collect() {
                    self.ack_cmd(msg.Cmd)
                } else {
                    self.nak_cmd(msg.Cmd)
                }
                return false
            }),
        link.Commands["Inventory"]:
            argnum(0, func(msg *link.Message) bool {
                self.send(link.NewMessage(link.Commands["Ack"], msg.Cmd, self.state.Inventory()))
                return false
            }),
        link.Commands["PrevResult"]:
            argnum(0, func(msg *link.Message) bool {
                self.send(link.NewMessage(link.Commands["Ack"], msg.Cmd, self.state.PrevResult()))
                return false
            }),
    }

    handle := func(msg *link.Message) bool {
        if f, ok := handlers[msg.Cmd]; ok {
            return f(msg)
        }
        panic(fmt.Sprintf("Command %s not found.", msg.Cmd))
    }

    complete := make(chan bool)
    self.state.NewMove()
    self.log.Println("Starting Turn", self.state.Turn)
    if !self.start_turn() {
        return false
    }
    go func(done chan<- bool) {
        for {
            if ok, msg := self.recv(); ok {
                if handle(msg) {
                    done <- true
                    break
                }
            } else {
                    done <- false
                    break
            }
        }
        return
    }(complete)
    c := <-complete
    self.log.Println("Ending Turn", self.state.Turn)
    return c
}

func (self *AgentProxy) start_turn() bool {
    return self.acked_send(link.NewMessage(link.Commands["Start"]))
}

func (self *AgentProxy) ack_cmd(cmd link.Command) {
    self.send(link.NewMessage(link.Commands["Ack"], cmd))
}

func (self *AgentProxy) nak_cmd(cmd link.Command) {
    self.send(link.NewMessage(link.Commands["Nak"], cmd))
}

func (self *AgentProxy) recv() (bool, *link.Message) {
    timeout := time.NewTicker(link.Timeout)
    select {
    case msg := <-self.rcv:
        self.log.Println("recv :", &msg)
        return true, &msg
    case <-timeout.C:
        timeout.Stop()
        self.log.Println("Client unresponsive.")
    }
    return false, nil
}

func (self *AgentProxy) send(msg *link.Message) bool {
    timeout := time.NewTicker(link.Timeout)
    select {
    case m := <-self.rcv:
        self.log.Println("recv unresolved message", m)
    case self.snd <- *msg:
        self.log.Println("sent :", msg)
        return true
    case <-timeout.C:
        timeout.Stop()
        self.log.Println("Client unresponsive.")
    }
    return false
}

func (self *AgentProxy) acked_send(msg *link.Message) bool {
    self.send(msg)
    return self.await_cmd_ack(msg.Cmd)
}

func (self *AgentProxy) await_cmd_ack(cmd link.Command) bool {
    if ok, msg := self.recv(); ok {
        if msg.Cmd == link.Commands["Ack"] && len(msg.Args) == 1 {
            switch acked := msg.Args[0].(type) {
            case link.Command:
                if acked == cmd {
                    return true
                }
            }
        }
    }
    return false
}
