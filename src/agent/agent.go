/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agent/agent.go
*/

package agent

import (
    "fmt"
    "json"
    "logflow"
    "net"
    "netchan"
    "util"
)
import (
    "agent/link"
    geo "coord/geometry"
    coordconf "coord/config"
)

type AgentComm chan []byte

type AgentConfig struct {
    Id uint32
    Position geo.Point
    Energy int
    Logs coordconf.LogConfigList
    CoordAddress string
}

type Agent interface {
    Turn(Comm)
    Id() uint32
    SetId(uint32)
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
                logger.Print("Starting")
                start()
            case msg.Cmd == link.Commands["Exit"]:
                logger.Print("Exiting")
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

func channelsForCoordinator(id uint32, addr string) (chan link.Message, chan link.Message) {
    ch_send := make(chan link.Message)
    ch_recv := make(chan link.Message)

    logflow.Print("agent", "Importing ", fmt.Sprintf("agent_req_%d", id))
    
    imp := util.MakeImporterWithRetry("tcp", addr, 10, logflow.NewSource("agent"))

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

func RunStandalone(myAddr string, agent Agent) {
    e := netchan.NewExporter()
    masterReq := make(AgentComm)
    masterRsp := make(AgentComm)
    e.Export("master_req", masterReq, netchan.Recv)
    e.Export("master_rsp", masterRsp, netchan.Send)
    
    log := logflow.NewSource("agent")
    
    log.Print("Listening at ", myAddr)
    
    addr, err := net.ResolveTCPAddr(myAddr)
    if err != nil {
        log.Fatal(err)
    }
    
    listener, err := net.ListenTCP(addr.Network(), addr)
    if err != nil {
        log.Fatal(err)
    }
    
    // RACE CONDITION!
    conn, err := listener.AcceptTCP()
    log.Print("Serving netchan master export")
    if err != nil {
        log.Fatal("listen:", err)
    }
    
    conn.SetLinger(0)
    go e.ServeConn(conn)
    
    log.Print("Closing listener")
    listener.Close()
    
    aconf := new(AgentConfig)
    
    bytes := <- masterReq
    
    err = json.Unmarshal(bytes, aconf)
    if err != nil {
        log.Fatal(err)
    }
    
    masterRsp <- []byte("configured")
    <- masterReq // "go"
    
    // Ignore position and energy
    agent.SetId(aconf.Id)
    aconf.Logs.Apply()
    RunWithCoordinator(agent, aconf.CoordAddress)
}
