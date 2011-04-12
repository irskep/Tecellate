package agent

import geo "coord/geometry"

type Message struct {
    msg []byte
    frequency uint8
    source *geo.Point
}

func (self *AgentState) NewMessage(freq uint8, msg []byte) *Message {
    m := &Message{frequency:freq, msg:msg, source:self.Position}
    self.Move.Messages = append(self.Move.Messages, m)
    return m
}

func (self *Message) Message() []byte { return self.msg }
func (self *Message) Frequency() uint8 { return self.frequency }
func (self *Message) Source() *geo.Point { return self.source }
