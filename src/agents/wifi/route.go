package wifi

import "fmt"
import . "byteslice"

type Route struct {
    Hops uint16
    DestAddr uint32
    NextAddr uint32
}

func NewRoute(hops uint16, dest_addr, next_addr uint32) *Route {
    return &Route{Hops:hops, DestAddr:dest_addr, NextAddr:next_addr}
}

func MakeRoute(from uint32, bytes ByteSlice) *Route {
    return &Route{
        Hops:bytes[:2].Int16(),
        DestAddr:bytes[2:].Int32(),
        NextAddr:from,
    }
}

func (self *Route) IncHops() {
    self.Hops += 1
}

func (self *Route) Bytes() ByteSlice {
    bytes := make(ByteSlice, 6)
    copy(bytes[:2], ByteSlice16(self.Hops))
    copy(bytes[2:], ByteSlice32(self.DestAddr))
    return bytes
}

func (self *Route) String() string {
    return fmt.Sprintf("<Route hops:%v dest:%v next:%v>", self.Hops, self.DestAddr, self.NextAddr)
}
