package game

import "crypto/rand"
import geo "coord/geometry"

const MessageLength = 256

type Message interface {
    Source() *geo.Point
    Message() []byte
    Frequency() uint8
}

type sortableMessages struct {
    msgs []Message
    targ *geo.Point
}

type Messages map[uint8][]Message

// convience functions ---------------------------------------------------------
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

// Messsages Methods -----------------------------------------------------------
func (self Messages) Add(msg Message) {
    f := msg.Frequency()
    if _, has := self[f]; !has {
        self[f] = make([]Message, 0, 10)
    }
    self[f] = append(self[f], msg)
}

func (self Messages) Hear(loc *geo.Point, freq uint8) []byte {
    if messages, has := self[freq]; has {
        msgs := newSortableMessages(len(messages), loc)
        for _, msg := range messages {
            msgs.add(msg)
        }
        return nil
    }
    return randbytes(MessageLength)
}

// messageSlice Methods --------------------------------------------------------
func newSortableMessages(size int, loc *geo.Point) *sortableMessages {
    return &sortableMessages{
            msgs:make([]Message, size),
            targ:loc,
    }
}
func (self *sortableMessages) add(msg Message) {
    self.msgs = append(self.msgs, msg)
}

// sort interface
func (self *sortableMessages) Len() int { return len(self.msgs) }
func (self *sortableMessages) Swap(i, j int) {
    self.msgs[i], self.msgs[j] = self.msgs[j], self.msgs[i]
}
func (self *sortableMessages) Less(i, j int) bool {
    a := self.msgs[i].Source().DistanceSquare(self.targ)
    b := self.msgs[j].Source().DistanceSquare(self.targ)
    return a < b
}
