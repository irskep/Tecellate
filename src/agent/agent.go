/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agent/agent.go
*/

package agent

import (
    "fmt"
    "logflow"
)
import (
    "agent/link"
)

type Agent interface {
    Turn(Comm)
    Id() uint
}

func Run(agent Agent, send link.SendLink, recv link.RecvLink) {
    logger := logflow.NewSource(fmt.Sprintf("agent/%d", agent.Id()))
    comm := StartComm(send, recv, logger)
    complete := make(chan bool)
    go func(send link.SendLink, recv link.RecvLink, done chan<- bool) {
        start := func() {
            comm.ack_start()
            agent.Turn(comm)
            comm.complete()
        }
        id := func() {
            comm.id(agent.Id())
        }

        for {
            switch msg := comm.recv_forever(); {
            case msg.Cmd == link.Commands["Start"]:
                start()
            case msg.Cmd == link.Commands["Id"]:
                id()
            case msg.Cmd == link.Commands["Exit"]:
                break
            default:
                s := fmt.Sprintf("Command %s not valid for current state.", msg.Cmd)
                panic(s)
            }
        }
        done <- true
    }(send, recv, complete)
    if ok := <-complete; ok {
        return
    }
    panic("we had an issue.")
}
