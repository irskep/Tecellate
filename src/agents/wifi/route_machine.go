package wifi

import "fmt"
import pseudo_rand "rand"
// import crypto_rand "crypto/rand"
import "agent"
import "logflow"
import . "byteslice"

const ROUTE_HOLDTIME = 25
const ROUTE_PAUSE = 10

type RouteMachine struct {
    agent agent.Agent
    logger logflow.Logger
    last ByteSlice
    state uint32
    backoff uint32
    wait uint32
    next_state uint32
    routes map[uint32]*Route
    route_keys []uint32
    next_route int
}

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

func NewRouteMachine(agent agent.Agent) *RouteMachine {
    self := &RouteMachine {
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/route/%d", agent.Id())),
        agent:agent,
        backoff:BACKOFF,
        wait:ROUTE_HOLDTIME,
        state:2,
        next_state:0,
        next_route:0,
        routes:make(map[uint32]*Route),
    }
    id := uint32(self.agent.Id())
    self.routes[id] = NewRoute(0, id, id)
    self.set_route_keys()
    return self
}

func (self *RouteMachine) Run(neighbors []uint32, comm agent.Comm) {
    for _, neighbor := range neighbors {
        self.routes[neighbor] = NewRoute(1, neighbor, neighbor)
    }
    self.set_route_keys()
    self.PerformListens(comm)
    if self.agent.Time()%10 == 9 {
        s := fmt.Sprintf("\nRoute Table (%v):\n", self.agent.Id())
        for i := uint32(1); i <= 8; i++ {
            if route, has := self.routes[i]; has {
                s += fmt.Sprint(route, "\n")
            }
        }
        self.log("info", s)
    }
    self.PerformSends(comm)
}

func (self *RouteMachine) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *RouteMachine) set_route_keys() {
    self.route_keys = make([]uint32, 0, len(self.routes))
    for k, _ := range self.routes {
        self.route_keys = append(self.route_keys, k)
    }
}

func (self *RouteMachine) confirm_last(comm agent.Comm) (confirm bool) {
    bytes := comm.Listen(2)
    confirm = self.last.Eq(bytes)
    self.log("info", self.agent.Time(), "confirm_last", confirm)
    return
}

func (self *RouteMachine) send_route(comm agent.Comm) {
    var route *Route
    if self.next_route < len(self.route_keys) {
        route = self.routes[self.route_keys[self.next_route]]
    } else {
        self.next_route = 0
        return
    }
    pkt := NewPacket(Commands["ROUTE"], uint32(self.agent.Id()))
    pkt.SetBody(route.Bytes())
    bytes := pkt.Bytes()
    comm.Broadcast(2, bytes)
    self.last = bytes
    self.log("info", self.agent.Time(), "sent", pkt, route)
}

func (self *RouteMachine) PerformSends(comm agent.Comm) {
    switch self.state {
        case 0:
            self.send_route(comm)
            self.state = 1
        case 1:
            self.next_state = 0
            if self.confirm_last(comm) {
                self.backoff = BACKOFF
                self.next_route += 1
                self.state = 2
                if self.next_route >= len(self.route_keys) {
                    self.next_route = 0
                    self.wait = ROUTE_HOLDTIME
                } else {
                    self.wait = ROUTE_PAUSE
                }
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

func (self *RouteMachine) PerformListens(comm agent.Comm) {
    switch self.state {
        case 1:
            return
    }
    pkt := MakePacket(comm.Listen(2))
    if !pkt.ValidateChecksum() { return }
    ok, cmd, _ := pkt.Cmd()
    if !ok { return }
    switch cmd {
        case Commands["ROUTE"]:
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

            self.log("info", self.agent.Time(), "Got a route from", from, "route", route)
    }
}
