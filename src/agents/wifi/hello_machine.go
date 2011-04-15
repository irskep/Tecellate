package wifi


import "fmt"
import pseudo_rand "rand"
// import crypto_rand "crypto/rand"
import "agent"
import "logflow"
import . "byteslice"


const (
    BACKOFF = 3
    HOLDTIME = 15
    RESET = 250
)

type HelloMachine struct {
    agent agent.Agent
    logger logflow.Logger
    last ByteSlice
    state uint32
    backoff uint32
    wait uint32
    next_state uint32
    neighbors map[uint32]uint32
}

func NewHelloMachine(agent agent.Agent) *HelloMachine {
    return &HelloMachine {
        agent:agent,
        backoff:BACKOFF,
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/hello/%d", agent.Id())),
        neighbors:make(map[uint32]uint32),
    }
}

func (self *HelloMachine) Run(comm agent.Comm) {
    self.PerformListens(comm)
    self.PerformSends(comm)
}

func (self *HelloMachine) Neighbors() []uint32 {
    neighbors := make([]uint32, 0, len(self.neighbors))
    for id, time := range self.neighbors {
        if time + RESET > uint32(self.agent.Time()) {
            neighbors = append(neighbors, id)
        }
    }
    return neighbors
}

func (self *HelloMachine) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *HelloMachine) confirm_last(comm agent.Comm) (confirm bool) {
    bytes := comm.Listen(1)
    confirm = self.last.Eq(bytes)
//     self.log("info", self.agent.Time(), "confirm_last", confirm)
    return
}

func (self *HelloMachine) hello(comm agent.Comm) {
    pkt := MakeHello(uint32(self.agent.Id()))
    bytes := pkt.Bytes()
//     self.log("info", self.agent.Time(), "sending", pkt)
    comm.Broadcast(1, bytes)
    self.last = bytes
}

func (self *HelloMachine) PerformSends(comm agent.Comm) {
    switch self.state {
        case 0:
            self.hello(comm)
            self.state = 1
        case 1:
            self.next_state = 0
            if self.confirm_last(comm) {
//                 self.log("info", self.agent.Time(), "restart")
                self.state = 2
                self.backoff = BACKOFF
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

func (self *HelloMachine) PerformListens(comm agent.Comm) {
    switch self.state {
        case 1:
            return
    }
    pkt := MakePacket(comm.Listen(1))
    if !pkt.ValidateChecksum() { return }
    ok, cmd, _ := pkt.Cmd()
    if !ok { return }
    switch cmd {
        case Commands["HELLO"]:
            id := pkt.IdField()
            self.neighbors[id] = uint32(self.agent.Time())
//             self.log("info", self.agent.Time(), "Got a hello from", id)
    }
}
