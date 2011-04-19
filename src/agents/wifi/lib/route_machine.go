package lib

import "fmt"
import pseudo_rand "rand"
import "agent"
import "logflow"
import . "byteslice"

import . "agents/wifi/lib/route"
import . "agents/wifi/lib/packet"

const ROUTE_HOLDTIME = 30
const ROUTE_PAUSE = 7

type RoutingTable map[uint32]*Route

type RouteMachine struct {
    freq uint8
    agent agent.Agent
    logger logflow.Logger
    last ByteSlice
    state uint32
    backoff float64
    wait uint32
    next_state uint32
    routes RoutingTable
    route_keys []uint32
    next_route int
    confirmed uint32
    not_confirmed uint32
}

func NewRouteMachine(freq uint8, agent agent.Agent) *RouteMachine {
    self := &RouteMachine {
        freq:freq,
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/route/%d", agent.Id())),
        agent:agent,
        backoff:BACKOFF,
        wait:uint32(float64(ROUTE_HOLDTIME)*(pseudo_rand.Float64() + 1.5)),
        state:2,
        next_state:0,
        next_route:0,
        routes:make(RoutingTable),
    }
    id := uint32(self.agent.Id())
    self.routes[id] = NewRoute(0, id, id)
    self.set_route_keys()
    return self
}

func (self *RouteMachine) Reachable() []uint32 {
    return self.route_keys
}

func (self *RouteMachine) Routes() RoutingTable {
    return self.routes
}

func (self *RouteMachine) Run(neighbors []uint32, comm agent.Comm) {
    for _, neighbor := range neighbors {
        self.routes[neighbor] = NewRoute(1, neighbor, neighbor)
    }
    self.clean_table()
    self.PerformListens(comm)
    self.PerformSends(comm)
}

func (self *RouteMachine) ConfirmRate() float64 {
    return float64(self.confirmed)/float64(self.confirmed+self.not_confirmed)
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

func (self *RouteMachine) clean_table() {
    self.set_route_keys()
    for _, k := range self.route_keys {
        route := self.routes[k]
        if route.DestAddr == uint32(self.agent.Id()) { continue }
        route.DecTTL()
        if route.TTL == 0 {
            self.routes[k] = nil, false
        }
    }
    self.set_route_keys()
}

func (self *RouteMachine) confirm_last(comm agent.Comm) (confirm bool) {
    bytes := comm.Listen(self.freq)
    confirm = self.last.Eq(bytes)
//     self.log("info", self.agent.Time(), "confirm_last", confirm)
    if confirm { self.confirmed += 1} else { self.not_confirmed += 1}
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
    pkt := NewPacket(Commands["ROUTE"], self.agent.Id(), 0xffffffff)
    pkt.SetBody(route.Bytes())
    bytes := pkt.Bytes()
    comm.Broadcast(self.freq, bytes)
    self.last = bytes
//     self.log("info", self.agent.Time(), "sent", pkt, route)
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
                    self.wait = uint32(float64(ROUTE_HOLDTIME)*(pseudo_rand.Float64() + 1.5))
                } else {
                    self.wait = uint32(float64(ROUTE_PAUSE)*(pseudo_rand.Float64() + 1.5))
                }
            } else {
                self.state = 2
                self.backoff = self.backoff*(pseudo_rand.Float64()*2 + 1)
                self.wait = uint32(self.backoff)
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
    pkt := MakePacket(comm.Listen(self.freq))
    if !pkt.ValidateChecksum() { return }
    ok, cmd, _ := pkt.Cmd()
    if !ok { return }
    switch cmd {
        case Commands["ROUTE"]:
            from := pkt.FromField()
            body := pkt.GetBody(6)
            route := MakeRoute(from, body)
            route.IncHops()

            if cur, has := self.routes[route.DestAddr]; has {
                if route.Hops < cur.Hops {
                    self.routes[route.DestAddr] = route
                } else if route.Hops == cur.Hops {
                    self.routes[route.DestAddr] = route
                }
            } else {
                self.routes[route.DestAddr] = route
            }

//             self.log("info", self.agent.Time(), "Got a route from", from, "route", route)
    }
}
