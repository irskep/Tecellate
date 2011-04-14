package wifi

import "testing"

import (
    "os"
)
import (
    "fmt"
    "agent"
    "agent/link"
    "coord"
    cagent "coord/agent"
    aproxy "coord/agent/proxy"
    geo "coord/geometry"
    "logflow"
)


func initLogs(name string, t *testing.T) func() {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs/wifi", 0776)
    if err != nil {
        panic("Directory logs/wifi could not be created.")
    }
    logflow.FileSink("logs/wifi/test/" + name, true, ".*")
    logflow.StdoutSink("agent/wifi.*")

    defer logflow.Println("test", fmt.Sprintf(`
--------------------------------------------------------------------------------
    Start Testing %v
`, name))
    return func() {
    logflow.Println("test", fmt.Sprintf(`
--------------------------------------------------------------------------------
    End Testing %v
`, name))
    logflow.RemoveAllSinks()
    }
}

func makeAgent(id uint, pos *geo.Point, energy cagent.Energy) *aproxy.AgentProxy {
    agnt := make(chan link.Message, 10)
    prox := make(chan link.Message, 10)
    simple := NewWifiBot(id)
    proxy := aproxy.NewAgentProxy(prox, agnt)
    proxy.SetState(cagent.NewAgentState(0, pos, energy))
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    return proxy
}

func TestAnnounce_2Agents(t *testing.T) {
    defer initLogs("TestAnnounce_2Agents", t)()

    gameconf := coord.NewGameConfig(10, "noise", true, false, 50, 50)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(1, 1), 10))
    gameconf.AddAgent(makeAgent(2, geo.NewPoint(3, 3), 10))
//     gameconf.AddAgent(makeAgent(3, geo.NewPoint(5, 5), 10))

    coords := gameconf.InitWithChainedLocalCoordinators(1, 50)
    coords.Run()
}
