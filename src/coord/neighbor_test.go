package coord

import "agents/configurable"
import "agent"

import (
    "fmt"
    "logflow"
    "os"
    "testing"
)

func initLogs(name string, t *testing.T) func() {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs/neighbor_test", 0776)
    if err != nil {
        panic("Directory logs/neighbor_test could not be created.")
    }

    logflow.FileSink("logs/neighbor_test/all", true, ".*")
    logflow.FileSink("logs/neighbor_test/" + name, true, ".*")
    logflow.FileSink("logs/neighbor_test/debug", true, ".*/debug")

    // Or show all output anyway I guess...
    // logflow.StdoutSink(".*")
//     logflow.StdoutSink(".*/debug")

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

func makeAgent(id uint32, xVelocity, yVelocity int) agent.Agent {
    a := configurable.New(id)
    a.XVelocity = xVelocity
    a.YVelocity = yVelocity
    a.LogMove = true
    return a
}

func TestLocalInfoPass(t *testing.T) {
    // initLogs("Local info", t)
    //
    // gameconf := NewGameConfig(11, "noise", false, true, 20, 10)
    // gameconf.AddAgent(makeAgentLocal(1, 0, 0))
    //
    // coords := gameconf.InitWithChainedLocalCoordinators(2, 10)
    // coords.Run()
    //
    // logflow.RemoveAllSinks()
}

func TestTCPInfoPass(t *testing.T) {
    defer initLogs("TCP info", t)()

    logflow.FileSink("logs/neighbor_test/agents", true, "test|agent/.*")
//     logflow.StdoutSink(".*/(info|debug)")

    gameconf := NewGameConfig(3, "noise", false, 20, 10)

    // gameconf.AddAgent(id, x, y)
    gameconf.AddAgent(1, 9, 0, 100)
    // gameconf.AddAgent(2, 5, 0, 100)

    agents := map[uint32]agent.Agent{1: makeAgent(1, 1, 0)}//, 2: makeAgent(2, -1, 0)}

    coords := gameconf.InitWithTCPChainedLocalCoordinators(2)
    // Start/connect the agents
    coords.StartAndConnectAgents(agents)
    // Run the coordinators
    coords.Run()
}
