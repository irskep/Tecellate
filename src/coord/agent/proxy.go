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
                    fmt.Println("proxy recieved a message")
                    fmt.Println(msg)
                    break loop
                case <-timeout.C:
                    fmt.Println("Timeout")
                    timeout.Stop()
                    break loop
            }
            println("loop")
        }
        println("end loop")
        done <- true
        return
    }(complete)
    c := <-complete
    println(c)
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
