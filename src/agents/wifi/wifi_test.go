package wifi

import "testing"

import (
    "os"
//     "runtime"
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

/*
func init() {
    runtime.GOMAXPROCS(2)
}*/

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
    proxy.SetState(cagent.NewAgentState(0, *pos, energy))
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    return proxy
}

func TestAnnounce_2Agents(t *testing.T) {
    defer initLogs("TestAnnounce_2Agents", t)()

    var time cagent.Energy = 700
    gameconf := coord.NewGameConfig(int(time), "noise", true, false, 50, 50)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(1, 1), time))
    gameconf.AddAgent(makeAgent(2, geo.NewPoint(3, 3), time))
    gameconf.AddAgent(makeAgent(3, geo.NewPoint(5, 5), time))
    gameconf.AddAgent(makeAgent(4, geo.NewPoint(7, 7), time))
    gameconf.AddAgent(makeAgent(5, geo.NewPoint(9, 9), time))
    gameconf.AddAgent(makeAgent(6, geo.NewPoint(11, 11), time))
    gameconf.AddAgent(makeAgent(7, geo.NewPoint(13, 13), time))
    gameconf.AddAgent(makeAgent(8, geo.NewPoint(15, 15), time))
    gameconf.AddAgent(makeAgent(9, geo.NewPoint(17, 17), time))

    coords := gameconf.InitWithChainedLocalCoordinators(1, 50)
    coords.Run()
}
