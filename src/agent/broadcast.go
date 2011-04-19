package agent

import "fmt"
import . "byteslice"
type Broadcast struct {
    freq uint8
    msg []byte
}

func newBroadcast(freq uint8, msg []byte) *Broadcast {
    self := new(Broadcast)
    self.freq = freq
    self.msg = msg
    return self
}


func MakeBroadcast(bytes ByteSlice) *Broadcast {
    freq  := bytes[0:1]
    msg   := bytes[1:]
    return newBroadcast(freq.Int8(), msg)
}

func (self *Broadcast) Message() (uint8, []byte) {
    return self.freq, self.msg
}

func (self *Broadcast) String() string {
    return fmt.Sprintf("on %d message '%s'", self.freq, self.msg)
}

func (self *Broadcast) Bytes() ByteSlice {
    bytes := make(ByteSlice, len(self.msg)+1)
    freq  := bytes[0:1]
    msg   := bytes[1:]
    copy(freq, ByteSlice8(self.freq))
    copy(msg, self.msg)
    return bytes
}
