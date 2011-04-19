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
    logger logflow.Logger
    time uint
    hello *HelloMachine
    route *RouteMachine
    send  *SendMachine
}

func NewStaticBot(id uint32) *StaticBot {
    self := &StaticBot{
        id:uint32(id),
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/static/%d", id)),
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

    if self.Id() == 8 && self.Time() == 500 {
        self.send.Send([]byte("Hello there Number 1."), 1)
    }

    self.hello.Run(comm)
    self.route.Run(self.hello.Neighbors(), comm)
    m := self.send.Run(self.route.Routes(), comm)


    if m != nil {
        self.log("info", self.Time(), "got a message", string([]byte(m.Body())))
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

