package agent

import "fmt"

type Broadcast struct {
    Freq uint8
    Msg []byte
}

func NewBroadcast(freq uint8, msg []byte) Broadcast {
    self := new(Broadcast)
    self.Freq = freq
    self.Msg = msg
    return *self
}

func (self Broadcast) Message() (uint8, []byte) {
    return self.Freq, self.Msg
}

func (self Broadcast) String() string {
    return fmt.Sprintf("on %d message '%s'", self.Freq, self.Msg)
}
