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

func initLogs(t *testing.T) {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")
    
    err := os.MkdirAll("logs", 0776)
    if err != nil {
        panic("Directory logs/ could not be created.")
    }
    
    logflow.FileSink("logs/NeighborTest_agents", "agent/.*")
    logflow.FileSink("logs/NeighborTest_agentproxies", "agentproxy/.*")
    logflow.FileSink("logs/NeighborTest_coords", "coord/.*")
    logflow.FileSink("logs/NeighborTest_coordproxies", "coordproxy/.*")
    logflow.FileSink("logs/NeighborTest_info", ".*/info")

    logflow.StdoutSink(".*")
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

func TestInfoPass(t *testing.T) {
    initLogs(t)
    
    logflow.Println("test", "\n\nTesting With 2 Coord and 2 Agents")
    logflow.Printf("main/info", "===New run: test info passing===")
    logflow.Printf("agent/all", "===New run: test info passing===")
    logflow.Printf("agentproxy/all", "===New run: test info passing===")
    logflow.Printf("coord/all", "===New run: test info passing===")
    logflow.Printf("coordproxy/all", "===New run: test info passing===")
    
    gameconf := NewGameConfig(11, "noise", false, true, 20, 10)
    gameconf.AddAgent(makeAgent(1, 0, 0))
    
    coords := gameconf.InitWithChainedLocalCoordinators(2, 10)
    coords.Run()
    
    logflow.RemoveAllSinks()
}
