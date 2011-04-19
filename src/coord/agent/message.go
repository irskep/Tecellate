package agent

import "fmt"
import geo "coord/geometry"

type Message struct {
    Msg []byte
    Frequency uint8
    Source geo.Point
}

func (self *AgentState) NewMessage(freq uint8, msg []byte) (bool, Message) {
    m := Message{Frequency:freq, Msg:msg, Source:self.Position}
    self.Move.Messages = append(self.Move.Messages, m)
    return true, m
}

func (self Message) String() string {
    return fmt.Sprintf("<cagent.Message source:%v freq:%v msg:\"%v\">", self.Source, self.Frequency, string(self.Msg))
}

