/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agent/agent.go
*/

package agent

import "fmt"
import "agent/link"

type Agent interface {
    Turn(Comm)
}

func Run(agent Agent, send, recv link.Link) {
    complete := make(chan bool)
    go func(send, recv link.Link, done chan<- bool) {
        start := func() {
            fmt.Println("Start Recieved")
            cm := StartComm(send, recv)
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
