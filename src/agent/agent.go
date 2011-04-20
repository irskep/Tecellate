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
    Id() uint32
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
            case msg.Cmd == link.Commands["Migrate"]:
                addr := string([]byte(msg.Args[0]))
                logflow.Print("agent/?/info", "I have been ordered to migrate to ", addr)
                comm.ack_migrate()
                ch_send, ch_recv := channelsForCoordinator(agent.Id(), addr)
                comm.SwapChannels(ch_send, ch_recv)
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
    var err os.Error
    for i := 0; i < 10; i++ {
        conn, err := net.Dial(network, "", remoteaddr)
        if err == nil {
            return netchan.NewImporter(conn)
        }
        logflow.Print("agent", "Netchan import failed, retrying")
        time.Sleep(1e9/2)
    }
    logflow.Print("agent", "Netchan import failed ten times. Bailing out.")
    logflow.Fatal("agent", err)
    return nil
}

func channelsForCoordinator(id uint32, addr string) (chan link.Message, chan link.Message) {
    ch_send := make(chan link.Message)
    ch_recv := make(chan link.Message)

    imp := makeImporterWithRetry("tcp", addr)

    logflow.Print("agent", "Importing ", fmt.Sprintf("agent_req_%d", id))

	err := imp.Import(fmt.Sprintf("agent_req_%d", id), ch_send, netchan.Send, 1)
	if err != nil {
	    logflow.Fatal("agent", err)
	}

    logflow.Print("agent", "Importing ", fmt.Sprintf("agent_rsp_%d", id))

	err = imp.Import(fmt.Sprintf("agent_rsp_%d", id), ch_recv, netchan.Recv, 1)
	if err != nil {
	    logflow.Fatal("agent", err)
	}
	return ch_send, ch_recv
}

func RunWithCoordinator(agent Agent, addr string) {
    ch_send, ch_recv := channelsForCoordinator(agent.Id(), addr)
	Run(agent, ch_send, ch_recv)
}
