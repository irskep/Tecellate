package simple

import "testing"

import (
    "os"
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

// func TestSimple(t *testing.T) {
//     fmt.Println("\n\nTesting Simple Turn Rollover")
//     agnt := make(chan link.Message, 10)
//     prox := make(chan link.Message, 10)
//     simple := NewSimple(1)
//     proxy := cagent.NewAgentProxy(prox, agnt)
//     go func() {
//         agent.Run(simple, agnt, prox)
//     }()
//     if !proxy.Turn() {
//         t.Error("Turn did not complete.")
//     }
// }
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

func initLogs(t *testing.T) {
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs", 0776)
    if err != nil {
        panic("Directory logs/ could not be created.")
    }

    logflow.FileSink("logs/TestWith2Coord_2Agents_agents", "agent/.*")
    logflow.FileSink("logs/TestWith2Coord_2Agents_agentproxies", "agentproxy/.*")
    logflow.FileSink("logs/TestWith2Coord_2Agents_coords", "coord/.*")
    logflow.FileSink("logs/TestWith2Coord_2Agents_coordproxies", "coordproxy/.*")
    logflow.FileSink("logs/TestWith2Coord_2Agents_info", ".*info")

    logflow.StdoutSink(".*info")
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

func TestWith2Coord_2Agents(t *testing.T) {
    initLogs(t)

    logflow.Println("test", "\n\nTesting With 2 Coord and 2 Agents")

    gameconf := coord.NewGameConfig("noise", true, true, 50, 50)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(0, 0), 1))
    gameconf.AddAgent(makeAgent(2, geo.NewPoint(10, 1), 1))
    gameconf.AddAgent(makeAgent(3, geo.NewPoint(20, 1), 0))
    gameconf.AddAgent(makeAgent(4, geo.NewPoint(25, 1), 1))
    gameconf.AddAgent(makeAgent(5, geo.NewPoint(30, 1), 2))

    coords := gameconf.InitWithChainedLocalCoordinators(1, 10)
    coords.Run()

    logflow.RemoveAllSinks()
}
