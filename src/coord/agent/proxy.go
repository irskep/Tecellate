package agent

import "fmt"
import "time"
import "agent/link"

type AgentProxy struct {
    State AgentState
    conn link.Link
}

func NewAgentProxy(conn link.Link) *AgentProxy {
    self := new(AgentProxy)
//     self.State = NewAgentState()
    self.conn = conn
    return self
}

func (self *AgentProxy) Turn() bool {
    var handlers = map[link.Command](func(*link.Message) bool) {
        link.Commands["Move"]:
            func(msg *link.Message) bool {
                self.nak_cmd(msg.Cmd)
                return false
            },
        link.Commands["Complete"]:
            func(msg *link.Message) bool {
                self.ack_cmd(msg.Cmd)
                return true
            },
    }

    handle := func(msg *link.Message) bool {
        fmt.Println("proxy recieved a message")
        fmt.Println(msg)
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
                case msg = <-self.conn:
                    if handle(&msg) {
                        timeout.Stop()
                        done <- true
                        break loop
                    }
                    timeout = time.NewTicker(link.Timeout)
                case <-timeout.C:
                    fmt.Println("Agent Proxy Timed Out")
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
    self.conn <- *link.NewMessage(link.Commands["Start"])
    return self.await_cmd_ack("Start")
}

func (self *AgentProxy) await_cmd_ack(cmd string) bool {
    timeout := time.NewTicker(link.Timeout)
    select {
    case msg := <-self.conn:
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
    self.conn <- *link.NewMessage(link.Commands["Ack"], cmd)
}

func (self *AgentProxy) nak_cmd(cmd link.Command) {
    self.conn <- *link.NewMessage(link.Commands["Nak"], cmd)
}
