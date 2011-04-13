package coord

import (
    "logflow"
    "os"
    "testing"
)

func initLogs(t *testing.T) {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")
    
    err := os.MkdirAll("logs", 0776)
    if err != nil {
        panic("Directory logs/ could not be created.")
    }
    
    logflow.FileSink("logs/NeighborTest_agents", "agent/.*")
    logflow.FileSink("logs/NeighborTest_agentproxies", "agentproxy/.*")
    logflow.FileSink("logs/NeighborTest_coords", "coord/.*")
    logflow.FileSink("logs/NeighborTest_coordproxies", "coordproxy/.*")
    logflow.FileSink("logs/NeighborTest_info", ".*/info")

    logflow.StdoutSink(".*")
}

func TestInfoPass(t *testing.T) {
    initLogs(t)
    
    logflow.Printf("main/info", "===New run: test info passing===")
    logflow.Printf("agent/all", "===New run: test info passing===")
    logflow.Printf("agentproxy/all", "===New run: test info passing===")
    logflow.Printf("coord/all", "===New run: test info passing===")
    logflow.Printf("coordproxy/all", "===New run: test info passing===")
    
    gameconf := coord.NewGameConfig("noise", false, true, 20, 10)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(0, 0), 1))
    
    coords := gameconf.InitWithChainedLocalCoordinators(1, 10)
    coords.Run()
    
    logflow.RemoveAllSinks()
}
