package simple

import "testing"
import "agent"

func TestSimple(t *testing.T) {
    simple :=  NewSimple()
    agent.Run(simple, nil)
}

