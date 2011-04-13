package coord

import "agent"
import "agent/link"
import "agents/configurable"
import cagent "coord/agent"
import aproxy "coord/agent/proxy"
import geo "coord/geometry"

import (
    "logflow"
    "os"
    "testing"
)

func initLogs(t *testing.T, withStdout bool, initMsg string) {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs", 0776)
    if err != nil {
        panic("Directory logs/ could not be created.")
    }

    logflow.FileSink("logs/NeighborTest_agents", true, "agent/.*")
    logflow.FileSink("logs/NeighborTest_agentproxies", true, "agentproxy/.*")
    logflow.FileSink("logs/NeighborTest_coords", true, "coord/.*")
    logflow.FileSink("logs/NeighborTest_coordproxies", true, "coordproxy/.*")
    logflow.FileSink("logs/NeighborTest_info", true, ".*/info")

    if withStdout {
        logflow.StdoutSink(".*/info")
    }

    if len(initMsg) > 0 {
        logflow.Printf("main/info", initMsg)
        logflow.Printf("agent/all", initMsg)
        logflow.Printf("agentproxy/all", initMsg)
        logflow.Printf("coord/all", initMsg)
        logflow.Printf("coordproxy/all", initMsg)
    }
}

func makeAgent(id uint, x int, y int) *aproxy.AgentProxy {
    p2a := make(chan link.Message, 10)
    a2p := make(chan link.Message, 10)
    a := configurable.New(id)
    a.XVelocity = 1
    proxy := aproxy.NewAgentProxy(p2a, a2p)
    proxy.SetState(cagent.NewAgentState(0, geo.NewPoint(x, y), 0))
    go func() {
        agent.Run(a, a2p, p2a)
    }()
    return proxy
}

func TestLocalInfoPass(t *testing.T) {
    initLogs(t, false, "===New run: test local info passing===")

    gameconf := NewGameConfig(11, "noise", false, true, 20, 10)
    gameconf.AddAgent(makeAgent(1, 0, 0))

    coords := gameconf.InitWithChainedLocalCoordinators(2, 10)
    coords.Run()

    logflow.RemoveAllSinks()
}

func TestTCPInfoPass(t *testing.T) {
    initLogs(t, true, "===New run: test TCP info passing===")

    logflow.RemoveAllSinks()
    gameconf := NewGameConfig(2, "noise", false, true, 20, 10)
    gameconf.AddAgent(makeAgent(1, 0, 0))

    //gameconf.InitWithTCPChainedLocalCoordinators(2, 10)
    //coords.Run()

    logflow.RemoveAllSinks()
}
