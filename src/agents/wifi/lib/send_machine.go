package lib

import "fmt"
import pseudo_rand "rand"
import "container/list"
import "agent"
import "logflow"
import . "byteslice"

import . "agents/wifi/lib/message"
import . "agents/wifi/lib/packet"

type SendMachine struct {
    freq uint8
    agent agent.Agent
    logger logflow.Logger
    last ByteSlice
    state uint32
    backoff uint32
    wait uint32
    next_state uint32
    routes RoutingTable
    recieve ByteSlice
    sendq *MessageQueue
}

type MessageQueue struct {
    list *list.List
}

func NewMessageQueue() *MessageQueue {
    return &MessageQueue{
        list:list.New(),
    }
}

func (self *MessageQueue) Len() int {
    return self.list.Len()
}

func (self *MessageQueue) Empty() bool {
    return self.list.Len() == 0
}

func (self *MessageQueue) Queue(m *Message) {
    self.list.PushBack(m)
}

func (self *MessageQueue) Dequeue() (*Message, bool) {
    front := self.list.Front()
    if front == nil { return nil, false }
    m := front.Value.(*Message)
    self.list.Remove(front)
    return m, true
}

func (self *MessageQueue) Clean() {
    for e := self.list.Front(); e != nil; {
        m := e.Value.(*Message)
        m.DecTTL()
        if m.SendTTL == 0 || m.TTL == 0 {
            next_e := e.Next()
            self.list.Remove(e)
            e = next_e
            if e == nil { break }
        } else {
            e = e.Next()
        }
    }
}

func NewSendMachine(freq uint8, agent agent.Agent) *SendMachine {
    self := &SendMachine {
        freq:freq,
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/send/%d", agent.Id())),
        agent:agent,
        backoff:BACKOFF,
        wait:ROUTE_HOLDTIME,
        state:2,
        next_state:0,
        sendq:NewMessageQueue(),
    }
    return self
}

func (self *SendMachine) Run(routes RoutingTable, comm agent.Comm) *Message {
    self.routes = routes
    m := self.PerformListens(comm)
    self.PerformSends(comm)
    self.sendq.Clean()
    return m
}

func (self *SendMachine) Send(msg ByteSlice, dest uint32) {
    self.sendq.Queue(NewMessage(msg, uint32(self.agent.Id()), dest))
}

func (self *SendMachine) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *SendMachine) confirm_last(comm agent.Comm) (confirm bool) {
    bytes := comm.Listen(self.freq)
    confirm = self.last.Eq(bytes)
//     self.log("info", self.agent.Time(), "confirm_last", confirm)
    return
}

func (self *SendMachine) send_message(comm agent.Comm) bool {
    find := func() (*Message, bool) {
        for i := 0; i < self.sendq.Len(); i++ {
            msg, ok := self.sendq.Dequeue()
            if !ok { break; }
            if _, has := self.routes[msg.DestAddr]; !has {
                self.sendq.Queue(msg)
            } else {
                return msg, true
            }
        }
        return nil, false
    }
    msg, found := find()
    if !found { return false }
    next := self.routes[msg.DestAddr].NextAddr
    pkt := NewPacket(Commands["MESSAGE"], next)
    pkt.SetBody(msg.Bytes())
    bytes := pkt.Bytes()
    comm.Broadcast(self.freq, bytes)
    self.last = bytes
    self.log("info", self.agent.Time(), "sent", pkt, msg)
    return true
}

func (self *SendMachine) PerformSends(comm agent.Comm) {
    switch self.state {
        case 0:
            if self.send_message(comm) {
                self.state = 1
            } else {
                self.state = 3
            }
        case 1:
            if self.confirm_last(comm) {
                self.state = 2
                self.next_state = 3
                self.backoff = BACKOFF
                self.wait = HOLDTIME
            } else {
                self.state = 2
                self.next_state = 0
                self.backoff = uint32(float64(self.backoff)*(pseudo_rand.Float64() + 1.5))
                self.wait = self.backoff
            }
        case 2:
            self.wait -= 1
            if self.wait == 0 {
                self.state = self.next_state
            }
        case 3:
            if !self.sendq.Empty() {
                self.state = 0
            }
        default:
//             self.log("debug", self.agent.Time(), "nop")
    }
}

func (self *SendMachine) PerformListens(comm agent.Comm) *Message {
    switch self.state {
        case 1:
            return nil
    }
    pkt := MakePacket(comm.Listen(self.freq))
    if !pkt.ValidateChecksum() { return nil }
    ok, cmd, _ := pkt.Cmd()
    if !ok { return nil }
    switch cmd {
        case Commands["MESSAGE"]:
            myaddr := uint32(self.agent.Id())
            to := pkt.IdField()
            body := pkt.GetBody(PacketBodySize)
            self.log("info", self.agent.Time(), "heard", to, "pkt", pkt, MakeMessage(body))
            if to == myaddr {
                msg := MakeMessage(body)
                if msg.ValidateChecksum() {
                    if msg.DestAddr == myaddr {
                        return msg
                    }
                    self.sendq.Queue(msg)
                }
            }
    }
    return nil
}
