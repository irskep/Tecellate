package simple

import "testing"

import "fmt"
import "agent"
import "agent/link"
import cagent "coord/agent"

func TestSimple(t *testing.T) {
    fmt.Println("Testing Simple Turn Rollover")
    agnt := make(link.Link, 10)
    prox := make(link.Link, 10)
    simple :=  NewSimple()
    proxy := cagent.NewAgentProxy(prox, agnt)
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    if !proxy.Turn() {
        t.Error("Turn did not complete.")
    }
}

