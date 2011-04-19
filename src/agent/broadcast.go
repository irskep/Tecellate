package agent

import "fmt"

type broadcast struct {
    freq uint8
    msg []byte
}

func newBroadcast(freq uint8, msg []byte) *broadcast {
    self := new(broadcast)
    self.freq = freq
    self.msg = msg
    return self
}

func (self *broadcast) Message() (uint8, []byte) {
    return self.freq, self.msg
}

func (self *broadcast) String() string {
    return fmt.Sprintf("on %d message '%s'", self.freq, self.msg)
}
