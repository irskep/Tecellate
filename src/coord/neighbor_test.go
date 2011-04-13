package simple

import "testing"

func initLogs(t *testing.T) {
    logflow.NewSink(logflow.NewTestWriter(t), ".*")
    
    err := os.MkdirAll("logs", 0776)
    if err != nil {
        panic("Directory logs/ could not be created.")
    }
    
    logflow.FileSink("logs/TestInfoPass_agents", "agent/.*")
    logflow.FileSink("logs/TestInfoPass_agentproxies", "agentproxy/.*")
    logflow.FileSink("logs/TestInfoPass_coords", "coord/.*")
    logflow.FileSink("logs/TestInfoPass_coordproxies", "coordproxy/.*")
    logflow.FileSink("logs/TestInfoPass_info", ".*info")

    logflow.StdoutSink(".*")
}

func TestInfoPass(t *testing.T) {
    initLogs(t)
    
    gameconf := coord.NewGameConfig("noise", false, true, 20, 10)
    gameconf.AddAgent(makeAgent(1, geo.NewPoint(0, 0), 1))
    
    coords := gameconf.InitWithChainedLocalCoordinators(1, 10)
    coords.Run()

    logflow.RemoveAllSinks()
}
