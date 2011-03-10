package simple

import "testing"

import "fmt"
import "agent"
import "agent/link"
import cagent "coord/agent"

func TestSimple(t *testing.T) {
    fmt.Println("Testing Simple Turn Rollover")
    conn := make(link.Link)
    simple :=  NewSimple()
    proxy := cagent.NewAgentProxy(conn)
    go func() {
        agent.Run(simple, conn)
    }()
    if !proxy.Turn() {
        t.Error("Turn did not complete.")
    }
}

