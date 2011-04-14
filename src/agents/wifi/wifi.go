/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agents/wifi/wifi.go
*/

package wifi

import "agent"
import "fmt"
import "logflow"

type WifiBot struct {
    id uint32
    logger logflow.Logger
    time uint32
    state uint32
}

func NewWifiBot(id uint) *WifiBot {
    self := &WifiBot{
        id:uint32(id),
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/%d", id)),
    }
    logflow.FileSink("logs/wifi/all", true, ".*")
    return self
}

func (self *WifiBot) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *WifiBot) Id() uint {
    return uint(self.id)
}

func (self *WifiBot) Turn(comm agent.Comm) {
    defer func(){self.time += 1}()
    self.hello(comm)
    pkt_1 := MakePacket(comm.Listen(1))
    self.log("info", self.time, pkt_1)
    return
}

func (self *WifiBot) hello(comm agent.Comm) {
    pkt := NewPacket(Commands["HELLO"])
    bytes := pkt.Bytes()
    comm.Broadcast(1, bytes)
}
