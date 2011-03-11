package agent

import (
//     "fmt"
    "time"
    "log"
    "os"
)
import (
    "agent/link"
    geo "coord/geometry"
)

type AgentProxy struct {
    State *AgentState
    snd link.Link
    rcv link.Link
    log *log.Logger
}

func NewAgentProxy(send, recv link.Link) *AgentProxy {
    self := new(AgentProxy)
    self.State = NewAgentState(0, geo.NewPoint(0, 0), 0)
    self.snd = send
    self.rcv = recv
    self.log = log.New(os.Stdout, "AgentProxy : ", 0)
    return self
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
            check_args(count, msg.Args)
            return f(msg)
        }
    }

    var handlers = map[link.Command]handler {
        link.Commands["Move"]:
            argnum(1, func(msg *link.Message) bool {
                mv := msg.Args[0].(link.Move).Move()
                if self.State.Move.Goto(&mv) {
                    self.ack_cmd(msg.Cmd)
                } else {
                    self.nak_cmd(msg.Cmd)
                }
                return false
            }),
        link.Commands["Complete"]:
            argnum(0, func(msg *link.Message) bool {
                self.ack_cmd(msg.Cmd)
                return true
            }),
    }

    handle := func(msg *link.Message) bool {
        self.log.Println("recv", msg)
        return handlers[msg.Cmd](msg)
    }

    complete := make(chan bool)
    if !self.start_turn() {
        return false
    }
    go func(done chan<- bool) {
        timeout := time.NewTicker(link.Timeout)
        loop: for {
            var msg link.Message
            select {
                case msg = <-self.rcv:
                    if handle(&msg) {
                        timeout.Stop()
                        done <- true
                        break loop
                    }
                    timeout = time.NewTicker(link.Timeout)
                case <-timeout.C:
                    self.log.Println("Agent Proxy Timed Out")
                    timeout.Stop()
                    done <- false
                    break loop
            }
        }
        return
    }(complete)
    c := <-complete
    return c
}

func (self *AgentProxy) start_turn() bool {
    self.snd <- *link.NewMessage(link.Commands["Start"])
    return self.await_cmd_ack("Start")
}

func (self *AgentProxy) await_cmd_ack(cmd string) bool {
    timeout := time.NewTicker(link.Timeout)
    select {
    case msg := <-self.rcv:
        self.log.Println(msg)
        timeout.Stop()
        if msg.Cmd == link.Commands["Ack"] && len(msg.Args) == 1 {
            switch acked := msg.Args[0].(type) {
            case link.Command:
                if acked == link.Commands[cmd] {
                    return true
                }
            }
        }
    case <-timeout.C:
        timeout.Stop()
    }
    return false
}

func (self *AgentProxy) ack_cmd(cmd link.Command) {
    self.snd <- *link.NewMessage(link.Commands["Ack"], cmd)
//     self.log.Println(link.NewMessage(link.Commands["Ack"], cmd))
}

func (self *AgentProxy) nak_cmd(cmd link.Command) {
    self.snd <- *link.NewMessage(link.Commands["Nak"], cmd)
}
