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
    sendq *MessageQueue
    lastrecv map[uint32]SEQUENCE
    send *SendMachine
    state uint32
    next_state uint32
}

func NewReliableSendMachine(send *SendMachine, agent agent.Agent) *ReliableSendMachine {
    self := &ReliableSendMachine {
        logger:logflow.NewSource(fmt.Sprintf("agent/wifi/Rsend/%d", agent.Id())),
        agent:agent,
        sendq:NewMessageQueue(),
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

    if !self.sendq.Empty() && !self.sendq.Peek().Sent {
        msg := self.sendq.Peek()
        self.send.Send(msg.Bytes(), msg.DestAddr)
        msg.Sent = false
    }

    if m != nil && !m.ValidateChecksum() { m = nil }
    if m != nil {
        msg := MakeMessage(m.Body())
        if msg.ValidateChecksum() {
            if msg.IsAck() && !self.sendq.Empty() {
                if msg.Acknowledge == self.sendq.Peek().Acknowledge {
                    self.sendq.Dequeue()
                }
            } else {
                self.lastrecv[m.FromAddr] = msg.Acknowledge
                ack := NewMessage(nil, self.lastrecv[m.DestAddr]+1, self.lastrecv[m.DestAddr], m.DestAddr)
                self.send.Send(ack.Bytes(), m.FromAddr)
                return msg
            }
        }
    }
    return nil
}

func (self *ReliableSendMachine) Send(msg ByteSlice, dest uint32) {
    m := NewMessage(msg, self.lastrecv[dest], self.lastrecv[dest] + 1, dest)
    self.lastrecv[dest] = self.lastrecv[dest] + 1
    self.sendq.Queue(m)
}
