package game

import geo "coord/geometry"

type Message interface {
    Source() *geo.Point
    Message() []byte
    Frequency() uint8
}

type Messages map[uint8][]Message

func (self Messages) Add(msg Message) {
    f := msg.Frequency()
    if _, has := self[f]; !has {
        self[f] = make([]Message, 0, 10)
    }
    self[f] = append(self[f], msg)
}
