package simple

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

func TestSimple(t *testing.T) {
    defer initLogs("TestSimple", t)()

    proxy := makeAgent(1, geo.NewPoint(0, 0), 1)
    gameconf := coord.NewGameConfig(3, "noise", true, false, 50, 50)
    gameconf.AddAgent(proxy)
    coord := gameconf.InitWithSingleLocalCoordinator()
    proxy.SetGameState(coord.GetGameState())
    if !proxy.Turn() {
        t.Error("Turn did not complete.")
    }
}
//
// func TestWithCoord(t *testing.T) {
//     fmt.Println("\n\nTesting With Coord")
//     agnt := make(chan link.Message, 10)
//     prox := make(chan link.Message, 10)
//     simple := NewSimple(1)
//     proxy := cagent.NewAgentProxy(prox, agnt)
//     go func() {
//         agent.Run(simple, agnt, prox)
//     }()
//
//     co := coord.NewCoordinator()
//     co.Configure(
//         config.NewConfig(
//             0,
//             []cagent.Agent{
//                 proxy,
//             },
//             "noise",
//             true,
//             true,
//             geo.NewPoint(0,0),
//             geo.NewPoint(10,10),
//         ),
//     )
//     co.Run()
// }

func initLogs(name string, t *testing.T) func() {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs/simple_test", 0776)
    if err != nil {
        panic("Directory logs/simple_test could not be created.")
    }

    logflow.FileSink("logs/simple_test/" + name, true, ".*")
    logflow.FileSink("logs/simple_test/all", true, ".*")

    // Or show all output anyway I guess...
//     logflow.StdoutSink(".*/info")
//     logflow.StdoutSink(".*/info")

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
    simple := NewSimple(id)
    proxy := aproxy.NewAgentProxy(prox, agnt)
    proxy.SetState(cagent.NewAgentState(0, pos, energy))
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    return proxy
}

func TestWithCoord_2Agents(t *testing.T) {
    defer initLogs("Coord_2Agents", t)()

    gameconf := coord.NewGameConfig(3, "noise", true, false, 50, 50)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(0, 0), 1))
//     gameconf.AddAgent(makeAgent(2, geo.NewPoint(10, 1), 1))
    gameconf.AddAgent(makeAgent(3, geo.NewPoint(20, 1), 1))
//     gameconf.AddAgent(makeAgent(4, geo.NewPoint(25, 1), 1))
    gameconf.AddAgent(makeAgent(5, geo.NewPoint(30, 1), 2))

    coords := gameconf.InitWithChainedLocalCoordinators(1, 50)
    coords.Run()


}
