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
    go func(done chan<- bool) {
        timeout := time.NewTicker(1e9) // timeout is 1 second
        loop: for {
            var msg link.Message
            select {
                case msg = <-self.conn:
                    fmt.Println(msg)
                case <-timeout.C:
                    fmt.Println("Timeout")
                    timeout.Stop()
                    break loop
            }
        }
        done <- true
    }(complete)
    self.start_turn()
    return <-complete
}

func (self *AgentProxy) start_turn() bool {
    self.conn <- *link.NewMessage(link.Commands["Start"])
    if msg := <- self.conn; msg.Cmd == link.Commands["Ack"] {
        return true
    }
    return false
}
