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
)

func TestSimple(t *testing.T) {
    fmt.Println("Testing Simple Turn Rollover")
    agnt := make(chan link.Message, 10)
    prox := make(chan link.Message, 10)
    simple := NewSimple(1)
    proxy := cagent.NewAgentProxy(prox, agnt)
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    if !proxy.Turn() {
        t.Error("Turn did not complete.")
    }
}

func TestWithCoord(t *testing.T) {
    fmt.Println("Testing With Coord")
    agnt := make(chan link.Message, 10)
    prox := make(chan link.Message, 10)
    simple := NewSimple(1)
    proxy := cagent.NewAgentProxy(prox, agnt)
    go func() {
        agent.Run(simple, agnt, prox)
    }()

    co := coord.NewCoordinator()
    co.Configure(
        config.NewConfig(
            0,
            []cagent.Agent{
                proxy,
            },
            "noise",
            true,
            true,
            geo.NewPoint(0,0),
            geo.NewPoint(10,10),
        ),
    )
//     fmt.Println("Testing Simple Turn Rollover")
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
}
