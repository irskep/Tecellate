package wifi

import "testing"

import (
    "os"
    "runtime"
    "fmt"
    "strings"
)
import (
    "agent"
    "agent/link"
    "coord"
    cagent "coord/agent"
    aproxy "coord/agent/proxy"
    geo "coord/geometry"
    "logflow"
)

var write_to_sinks logflow.WriteToSinks = logflow.WriteToSinksFunction
func init() {
    runtime.GOMAXPROCS(2)
}

func initLogs(name string, t *testing.T) func() {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs/wifi", 0776)
    if err != nil {
        panic("Directory logs/wifi could not be created.")
    }
//     logflow.FileSink("logs/wifi/test/" + name, true, ".*")
//     logflow.StdoutSink("agent/wifi.*")

    defer func() {
       logflow.WriteToSinksFunction = func(keypath, s string) {
           if strings.HasPrefix(keypath, "agent/wifi") {
               fmt.Print(keypath, ": ", s)
           }
       }
    }()

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
    logflow.WriteToSinksFunction = write_to_sinks
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

func TestAnnounce(t *testing.T) {
    defer initLogs("TestAnnounce", t)()

    var time cagent.Energy = 10000
    gameconf := coord.NewGameConfig(int(time), "noise", true, false, 100, 100)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(0, 0), time))
    gameconf.AddAgent(makeAgent(2, geo.NewPoint(6, 6), time))
    gameconf.AddAgent(makeAgent(3, geo.NewPoint(12, 12), time))
    gameconf.AddAgent(makeAgent(4, geo.NewPoint(18, 18), time))
    gameconf.AddAgent(makeAgent(5, geo.NewPoint(24, 24), time))
    gameconf.AddAgent(makeAgent(6, geo.NewPoint(30, 30), time))
    gameconf.AddAgent(makeAgent(7, geo.NewPoint(36, 36), time))
    gameconf.AddAgent(makeAgent(8, geo.NewPoint(42, 42), time))
    coords := gameconf.InitWithChainedLocalCoordinators(1, 60)
    coords.Run()
}
