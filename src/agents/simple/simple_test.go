package simple

import "testing"

import (
    "fmt"
)
import (
    "agent"
    "agent/link"
    "coord"
    "coord/config"
    cagent "coord/agent"
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
    logflow.FileSink("logs/TestWith2Coord_2Agents_agents", "agent/.*")
    
    // Proxies don't use numbers in their keypaths so don't show the prefixes
    // because they will all be identical
    ap, _ := logflow.FileSink("logs/TestWith2Coord_2Agents_agentproxies", "agentproxy/.*")
    ap.SetWritesPrefix(false)
    
    logflow.FileSink("logs/TestWith2Coord_2Agents_coords", "coord/.*")
    
    logflow.FileSink("logs/TestWith2Coord_2Agents_coords", ".*/info")
}

func makeAgent(id uint, pos *geo.Point, energy cagent.Energy) *cagent.AgentProxy {
    agnt := make(chan link.Message, 10)
    prox := make(chan link.Message, 10)
    simple := NewSimple(id)
    proxy := cagent.NewAgentProxy(prox, agnt)
    proxy.SetState(cagent.NewAgentState(0, pos, energy))
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    return proxy
}

func makeCoord(id int, tl, br *geo.Point, proxies []cagent.Agent) *coord.Coordinator {
    co := coord.NewCoordinator()
    co.Configure(
        config.NewConfig(
            id,
            proxies,
            "noise",
            true,
            true,
            tl,
            br,
        ),
    )
    return co
}

func TestWith2Coord_2Agents(t *testing.T) {
    fmt.Println("\n\nTesting With 2 Coord and 2 Agents")
    
    initLogs(t)
    
    proxies1 := make([]cagent.Agent, 0, 10)
//     proxies2 := make([]cagent.Agent, 0, 10)
    proxies1 = append(proxies1, makeAgent(1, geo.NewPoint(0, 0), 1))
    proxies1 = append(proxies1, makeAgent(2, geo.NewPoint(0, 1), 1))
    proxies1 = append(proxies1, makeAgent(3, geo.NewPoint(1, 1), 0))
    proxies1 = append(proxies1, makeAgent(4, geo.NewPoint(2, 1), 1))
    proxies1 = append(proxies1, makeAgent(5, geo.NewPoint(3, 1), 2))
//     proxies2 = append(proxies2, makeAgent(2))

//     fmt.Println(proxies)
//     for _, prox := range proxies {
//         fmt.Println(prox)
//     }
    coords := make(coord.CoordinatorSlice, 0, 10)
    coords = append(coords, makeCoord(1, geo.NewPoint(0,0),geo.NewPoint(9,9), proxies1))
//     coords = append(coords, makeCoord(2, geo.NewPoint(0,0),geo.NewPoint(9,9), proxies2))
    coord.ConnectInChain(coords)
    coords.Run()
    
    logflow.RemoveAllSinks()
}
