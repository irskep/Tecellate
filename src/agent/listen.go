package agent

import "fmt"
import . "byteslice"

type Listen struct {
    freq uint8
}

func newListen(freq uint8) *Listen {
    self := new(Listen)
    self.freq = freq
    return self
}

func MakeListen(bytes ByteSlice) *Listen {
    return newListen(bytes.Int8())
}

func (self *Listen) Listen() uint8 {
    return self.freq
}

func (self *Listen) String() string {
    return fmt.Sprintf("on %d", self.freq)
}

func (self *Listen) Bytes() ByteSlice {
    return ByteSlice8(self.freq)
}
