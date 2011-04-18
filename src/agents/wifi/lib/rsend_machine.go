package lib

import "fmt"
// import pseudo_rand "rand"
// import "container/list"
import "agent"
import "logflow"
import . "byteslice"

import . "agents/wifi/lib/message"
// import . "agents/wifi/lib/packet"

type ReliableSendMachine struct {
    agent agent.Agent
    logger logflow.Logger
    sent_queue map[uint32]*MessageQueue
    lastrecv map[uint32]SEQUENCE
    send *SendMachine
    state uint32
    next_state uint32
}

func NewReliableSendMachine(send *SendMachine, agent agent.Agent) *ReliableSendMachine {
    self := &ReliableSendMachine {
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/Rsend/%d", agent.Id())),
        agent:agent,
        sent_queue:make(map[uint32]*MessageQueue),
        lastrecv:make(map[uint32]SEQUENCE),
        send:send,
    }
    return self
}

func (self *ReliableSendMachine) log(level logflow.LogLevel, v ...interface{}) {
    self.logger.Logln(level, v...)
}

func (self *ReliableSendMachine) Run(routes RoutingTable, comm agent.Comm) *Message {
    m := self.send.Run(routes, comm)
    if m != nil && !m.ValidateChecksum() { m = nil }
    if m != nil {
        msg := MakeMessage(m.Body())
        if msg.ValidateChecksum() {
            if msg.IsAck() {
                if msg.Acknowledge == self.lastrecv[m.FromAddr] {
                    self.sent_queue[m.FromAddr].Dequeue()
                }
            } else {
                self.lastrecv[m.FromAddr] = msg.Acknowledge
                ack := NewMessage(nil, self.lastrecv[m.DestAddr]+1, self.lastrecv[m.DestAddr])
                self.send.Send(ack.Bytes(), m.FromAddr)
                return msg
            }
        }
    }
    return nil
}

func (self *ReliableSendMachine) Send(msg ByteSlice, dest uint32) {
    m := NewMessage(msg, self.lastrecv[dest], self.lastrecv[dest] + 1)
    self.lastrecv[dest] = self.lastrecv[dest] + 1
    if _, has := self.sent_queue[dest]; !has {
        self.sent_queue[dest] = NewMessageQueue()
    }
    self.sent_queue[dest].Queue(m)
    self.send.Send(msg, dest)
}
