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
    "net"
    "netchan"
    "os"
    "time"
)
import (
    "agent/link"
)

type Agent interface {
    Turn(Comm)
    Id() uint
    Time() uint
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

        for {
            switch msg := comm.recv_forever(); {
            case msg.Cmd == link.Commands["Start"]:
                start()
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

func makeImporterWithRetry(network string, remoteaddr string) *netchan.Importer {
    // This method is actually entirely futile because the race condition we're trying
    // to account for happens between listener creation and exporter.ServeConn().
    // An error is only thrown if the listener does not exist, but we must already
    // have a listener to call ServeConn().
    // To really fix this, you have to try sending a message down the pipe and see
    // if it panics.
    var err os.Error
    for i := 0; i < 3; i++ {
        conn, err := net.Dial(network, "", remoteaddr)
        if err == nil {
            return netchan.NewImporter(conn)
        }
        logflow.Print("agent", "Netchan import failed, retrying")
        time.Sleep(1e9/2)
    }
    logflow.Print("agent", "Netchan import failed three times. Bailing out.")
    logflow.Fatal("agent", err)
    return nil
}

func RunWithCoordinator(agent Agent, addr string) {
    ch_send := make(chan link.Message)
    ch_recv := make(chan link.Message)

    imp := makeImporterWithRetry("tcp", addr)

    logflow.Print("agent", "Importing ", fmt.Sprintf("agent_req_%d", agent.Id()))

	err := imp.Import(fmt.Sprintf("agent_req_%d", agent.Id()), ch_send, netchan.Send, 1)
	if err != nil {
	    logflow.Fatal("agent", err)
	}

    logflow.Print("agent", "Importing ", fmt.Sprintf("agent_rsp_%d", agent.Id()))

	err = imp.Import(fmt.Sprintf("agent_rsp_%d", agent.Id()), ch_recv, netchan.Recv, 1)
	if err != nil {
	    logflow.Fatal("agent", err)
	}

	Run(agent, ch_send, ch_recv)
}
