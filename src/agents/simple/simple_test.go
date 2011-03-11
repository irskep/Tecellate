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
    co := coord.NewCoordinator()
    co.Configure(config.NewConfig())
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
