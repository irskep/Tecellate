package game

import "fmt"
import pseudo_rand "rand"
import crypto_rand "crypto/rand"
import "sort"
import geo "coord/geometry"
import cagent "coord/agent"
import "logflow"

const MessageLength = 11
const hearing_range = 10.0
const corrupt_scale = 1.137
const combine_scale = corrupt_scale*3

var log logflow.Logger = logflow.NewSource(fmt.Sprint("message"))

type sortableMessages struct {
    msgs []cagent.Message
    targ *geo.Point
}
type Messages struct {
    Msgs map[uint8][]cagent.Message
    Cache map[complex128](map[uint8][]byte)
}

// initializer for random number generator -------------------------------------
func init() {
    // This function was a originally part of the structure/block/byteslice pkg
    // in SourceQL
    Int64 := func (b []byte) int64 {
        i := uint64(0)
        for j := 0; j < len(b) && j < 8; j++ {
            i |= 0x00000000000000ff & uint64(b[j])
            if j+1 < len(b) {
                i <<= 8
            }
        }
        return int64(i)
    }
    pseudo_rand.Seed(Int64(randbytes(8)))
}

// convience functions ---------------------------------------------------------
func randbyte() byte {
    return randbytes(1)[0]
}
func randbytes(k int) []byte {
    bytes := make([]byte, k)
    if n, err := crypto_rand.Read(bytes); n == k && err == nil {
        return bytes
    }
    panic("Can't get random byte.")
}
func corrupt(msg []byte, dist float64) (corrupted []byte) {
    corrupted = make([]byte, MessageLength)
    for i := 0; i < MessageLength; i++ {
        var cur byte
        if i < len(msg) { cur = msg[i] } else { cur = randbyte() }
        randfloat := pseudo_rand.Float64()
        prob := 1.0/(dist/(hearing_range*corrupt_scale))
        log.Logln(logflow.DEBUG, cur,
            dist > hearing_range && randfloat > prob,
            randfloat, prob)
        if dist > hearing_range && randfloat > prob {
            corrupted[i] = cur ^ randbyte()
//             log.Logln(logflow.DEBUG, "corrupted byte ", cur, corrupted[i])
        } else {
            corrupted[i] = cur
        }
    }
    return
}

// Messages Methods -----------------------------------------------------------
func NewMessages(peers []*GameStateResponse) *Messages {
    self := new(Messages)
    self.Msgs = make(map[uint8][]cagent.Message)
    self.Cache = make(map[complex128](map[uint8][]byte))
    for _, peer := range peers {
        for freq, msgs := range peer.Messages {
            self.Msgs[freq] = make([]cagent.Message, 0, len(msgs)+2)
            self.Msgs[freq] = append(self.Msgs[freq], msgs...)
        }
    }
    return self
}

func (self *Messages) Add(msg cagent.Message) {
    f := msg.Frequency
    if _, has := self.Msgs[f]; !has {
        self.Msgs[f] = make([]cagent.Message, 0, 10)
    }
    self.Msgs[f] = append(self.Msgs[f], msg)
}

func (self *Messages) Hear(loc *geo.Point, freq uint8) (msg []byte) {
    if freqs, has := self.Cache[loc.Complex()]; has {
        if m, has := freqs[freq]; has{
            log.Logln(logflow.DEBUG, "Cached!", loc, freq, m)
            return m
        }
    }

    msg = make([]byte, MessageLength)
    if messages, has := self.Msgs[freq]; has {
        log.Logln(logflow.DEBUG, "have a message on freq ", freq)
        msgs := newSortableMessages(len(messages), loc)
        for _, msg := range messages {
            msgs.add(msg)
        }
        msgs.sort()
        for i, M := range msgs.msgs {
            dist := M.Source.Distance(loc)
            log.Logln(logflow.DEBUG, "message", i, "dist to targ", dist)
            m := corrupt(M.Msg, dist)
            log.Logln(logflow.DEBUG, "message", i, "corrupted", string(m), m)
            if i == 0 {
                msg = m
            } else {
                for j, byt := range m {
                    randfloat := pseudo_rand.Float64()
                    prob := 1.0/(dist/(hearing_range/combine_scale))
                    decision := dist <= hearing_range || randfloat <= prob
                    log.Logln(logflow.DEBUG, msg[j], byt,
                        decision,
                        randfloat, prob)
                    if decision {
                        var r byte
                        if msg[j] == byt { r = msg[j] } else { r = msg[j] & (^byt) }
                        log.Logln(logflow.DEBUG, "combining", j, msg[j], byt, r)
                        msg[j] = r
                    }
                }
            }
            log.Logln(logflow.DEBUG, "message", i, "acc", string(msg), msg)
        }
    } else {
        log.Logln(logflow.DEBUG, "don't have a message on freq ", freq)
        msg = randbytes(MessageLength)
    }


    if _, has := self.Cache[loc.Complex()]; !has {
        self.Cache[loc.Complex()] = make(map[uint8][]byte)
    }
    self.Cache[loc.Complex()][freq] = msg

    return
}

// messageSlice Methods --------------------------------------------------------
func newSortableMessages(size int, loc *geo.Point) *sortableMessages {
    return &sortableMessages{
            msgs:make([]cagent.Message, 0, size),
            targ:loc,
    }
}
func (self *sortableMessages) add(msg cagent.Message) {
    self.msgs = append(self.msgs, msg)
}
func (self *sortableMessages) sort() *sortableMessages {
    sort.Sort(self)
    return self
}

// sort interface
func (self *sortableMessages) Len() int { return len(self.msgs) }
func (self *sortableMessages) Swap(i, j int) {
    self.msgs[i], self.msgs[j] = self.msgs[j], self.msgs[i]
}
func (self *sortableMessages) Less(i, j int) bool {
    a := self.msgs[i].Source.DistanceSquare(self.targ)
    b := self.msgs[j].Source.DistanceSquare(self.targ)
    return a < b
}
