/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agents/wifi/wifi.go
*/

package wifi

import "fmt"
import pseudo_rand "rand"
import crypto_rand "crypto/rand"
import "agent"
import "logflow"
import . "byteslice"

// initializer for random number generator -------------------------------------
func randbytes(k int) ByteSlice {
    bytes := make(ByteSlice, k)
    cbytes := bytes[:]
    for
        n, err := crypto_rand.Read(cbytes);
        n < k;
        n, err = crypto_rand.Read(cbytes) {
            if err != nil {
                panic("Can't get random bytes.")
            }
            k = k-n
            cbytes = cbytes[n:]
    }
    return bytes
}
func init() {
    pseudo_rand.Seed(int64(randbytes(8).Int64()))
}

// WifiBot ---------------------------------------------------------------------
type WifiBot struct {
    id uint32
    logger logflow.Logger
    time uint
    hello *HelloMachine
    route *RouteMachine
    send  *SendMachine
}

func NewWifiBot(id uint) *WifiBot {
    self := &WifiBot{
        id:uint32(id),
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/%d", id)),
    }
    self.hello = NewHelloMachine(1, self)
    self.route = NewRouteMachine(2, self)
    self.send = NewSendMachine(3, self)
//     logflow.FileSink("logs/wifi/all", true, ".*")
    return self
}

func (self *WifiBot) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *WifiBot) Time() uint {
    return self.time
}

func (self *WifiBot) Id() uint {
    return uint(self.id)
}

func (self *WifiBot) Turn(comm agent.Comm) {
    defer func(){self.time += 1}()

    if self.Id() == 8 && self.Time() == 200 {
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

