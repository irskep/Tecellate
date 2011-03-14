package agent

import "fmt"

type listen struct {
    freq uint8
}

func newListen(freq uint8) *listen {
    self := new(listen)
    self.freq = freq
    return self
}

func (self *listen) Listen() uint8 {
    return self.freq
}

func (self *listen) String() string {
    return fmt.Sprintf("on %d", self.freq)
}
