/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agent/agent.go
*/

package agent

import (
    "os"
    "fmt"
    "log"
)
import (
    "agent/link"
)

type Agent interface {
    Turn(Comm)
    Id() uint
}

func Run(agent Agent, send link.SendLink, recv link.RecvLink) {
    logger := log.New(os.Stdout, fmt.Sprintf("Agent(%d) : ", agent.Id()), 0)
    complete := make(chan bool)
    go func(send link.SendLink, recv link.RecvLink, done chan<- bool) {
        start := func() {
            logger.Println("Start Recieved")
            cm := StartComm(send, recv, logger)
            cm.ack_start()
            agent.Turn(cm)
            cm.complete()
        }

        loop: for {
            switch msg := <-recv; {
                case msg.Cmd == link.Commands["Start"]:
                    start()
                case msg.Cmd == link.Commands["Exit"]:
                    break loop
                default:
                    panic(
                        fmt.Sprintf(
                            "Command %s not valid for current state.",
                            msg.Cmd,
                        ),
                    )
            }
        }
        done <- true
    }(send, recv, complete)
    if ok := <-complete; ok {
        return
    }
    panic("we had an issue.")
}
