package wifi

import "fmt"
import pseudo_rand "rand"
// import crypto_rand "crypto/rand"
import "agent"
import "logflow"
import . "byteslice"

const MESSAGE_HOLD = 10

type SendMachine struct {
    agent agent.Agent
    logger logflow.Logger
    last ByteSlice
    state uint32
    backoff uint32
    wait uint32
    next_state uint32
    routes RoutingTable
    recieve ByteSlice
    sendqueue SendQueue
}

type SendQueue []*Message



func NewSendMachine(agent agent.Agent) *SendMachine {
    self := &SendMachine {
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/route/%d", agent.Id())),
        agent:agent,
        backoff:BACKOFF,
        wait:ROUTE_HOLDTIME,
        state:2,
        next_state:0,
    }
    return self
}

func (self *SendMachine) Run(routes RoutingTable, comm agent.Comm) {
    self.routes = routes
    self.PerformListens(comm)
    self.PerformSends(comm)
}

func (self *SendMachine) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *SendMachine) confirm_last(comm agent.Comm) (confirm bool) {
    bytes := comm.Listen(2)
    confirm = self.last.Eq(bytes)
//     self.log("info", self.agent.Time(), "confirm_last", confirm)
    return
}

func (self *SendMachine) send_message(comm agent.Comm) {
    pkt := NewPacket(Commands["MESSAGE"], uint32(self.agent.Id()))
//     pkt.SetBody(route.Bytes())
    bytes := pkt.Bytes()
    comm.Broadcast(2, bytes)
    self.last = bytes
}

func (self *SendMachine) PerformSends(comm agent.Comm) {
    switch self.state {
        case 0:
            self.send_message(comm)
            self.state = 1
        case 1:
            self.next_state = 0
            if self.confirm_last(comm) {
                self.backoff = BACKOFF
                self.state = 2
                self.wait = HOLDTIME
            } else {
                self.state = 2
                self.backoff = uint32(float64(self.backoff)*(pseudo_rand.Float64() + 1.5))
                self.wait = self.backoff
            }
        case 2:
//             self.log("debug", self.agent.Time(), "wait", self.wait, "backoff", self.backoff)
            self.wait -= 1
            if self.wait == 0 {
                self.state = self.next_state
            }
        case 3:
            fallthrough
        default:
//             self.log("debug", self.agent.Time(), "nop")
    }
}

func (self *SendMachine) PerformListens(comm agent.Comm) {
    switch self.state {
        case 1:
            return
    }
    pkt := MakePacket(comm.Listen(2))
    if !pkt.ValidateChecksum() { return }
    ok, cmd, _ := pkt.Cmd()
    if !ok { return }
    switch cmd {
        case Commands["MESSAGE"]:
            from := pkt.IdField()
            body := pkt.GetBody(6)
            route := MakeRoute(from, body)
            route.IncHops()

            if cur, has := self.routes[route.DestAddr]; has {
                if route.Hops < cur.Hops {
                    self.routes[route.DestAddr] = route
                }
            } else {
                self.routes[route.DestAddr] = route
            }

//             self.log("info", self.agent.Time(), "Got a route from", from, "route", route)
    }
}
