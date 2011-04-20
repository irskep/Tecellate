/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agents/wifi/wifi.go
*/

package wifi

import "fmt"
import "agent"
import "logflow"

import . "agents/wifi/lib"


// StaticBot ---------------------------------------------------------------------
type StaticBot struct {
    id uint32
    first uint32
    last uint32
    next uint32
    logger logflow.Logger
    time uint
    hello *HelloMachine
    route *RouteMachine
    send  *SendMachine
    recieved []uint32
}

func NewStaticBot(id, first, last uint32) *StaticBot {
    self := &StaticBot{
        id:id,
        first:first,
        last:last,
        next:first-1,
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/static/%d", id)),
        recieved:make([]uint32, 0, int(last-first)),
    }
    self.hello = NewHelloMachine(1, self)
    self.route = NewRouteMachine(15, self)
    self.send = NewSendMachine(3, self)
//     logflow.FileSink("logs/wifi/all", true, ".*")
    return self
}

func (self *StaticBot) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *StaticBot) Time() uint {
    return self.time
}

func (self *StaticBot) Id() uint32 {
    return self.id
}

func (self *StaticBot) Turn(comm agent.Comm) {
    defer func(){self.time += 1}()

//     self.log("Time = ", self.time)

    if self.Time()/uint(self.Id()) > 750 && self.next < self.last {
        self.next += 1
        if self.next == self.Id() { return }
        self.send.Send([]byte(fmt.Sprintf("Hello there Number %v.", self.next)), self.next)
//         self.next = self.last
    }

    self.hello.Run(comm)
    self.route.Run(self.hello.Neighbors(), comm)
    m := self.send.Run(self.route.Routes(), comm)


    if m != nil {
        self.log("info", self.Time(), "got a message", string([]byte(m.Body())))
//         self.send.Send([]byte(fmt.Sprintf("Thanks for the message %v", m.FromAddr)), m.FromAddr)
        self.recieved = append(self.recieved, m.FromAddr)
    }
    if self.Time()%100 == 9 {
//         self.log("info", self.Time(), "neighbors", self.hello.Neighbors())

//         self.log("info", self.Time(), "reachable", self.route.Reachable())
//         s := fmt.Sprintf("\nRoute Table (%v):\n", self.agent.Id())
//         for i := uint32(1); i <= 8; i++ {
//             if route, has := self.routes[i]; has {
//                 s += fmt.Sprint(route, "\n")
//             }
//         }
//         self.log("info", s)
    }
}

