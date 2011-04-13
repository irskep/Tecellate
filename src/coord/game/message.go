package game

import "crypto/rand"
import geo "coord/geometry"

const MessageLength = 256

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

func randbyte() byte {
    return randbytes(1)[0]
}

func randbytes(k int) []byte {
    bytes := make([]byte, k)
    if n, err := rand.Read(bytes); n == k && err == nil {
        return bytes
    }
    panic("Can't get random byte.")
}

func (self Messages) Hear(loc *geo.Point, freq uint8) []byte {
    if _, has := self[freq]; has {
        return nil
    }
    return randbytes(MessageLength)
}
