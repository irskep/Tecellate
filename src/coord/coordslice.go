package coord

import (
    agent "agent"
    aproxy "coord/agent/proxy"
    "logflow"
)

type CoordinatorSlice []*Coordinator

func (self CoordinatorSlice) Run() {
    // This channel will receive one 'true' for each process completion
    complete := make(chan bool)

    // Start the necessary threads
    for _, c := range(self) {
        c.StartRPCServer()
        go c.ProcessTurns(complete)
    }

    // Wait for processing to complete
    for _, _ = range(self) {
        <- complete
    }
    for _, c := range(self) {
        c.Close()
    }
}

func (self CoordinatorSlice) ConnectToLocalAgents(agents map[uint32]agent.Agent) {
    for _, c := range(self) {
        for _, ad := range(c.conf.Agents) {
            p := aproxy.RunAgentLocal(agents[ad.Id], ad.X, ad.Y)
            c.availableGameState.Agents = append(c.availableGameState.Agents, p)
        }
    }
}

func (self CoordinatorSlice) Chain() {
    for i, c := range(self) {
        if i < len(self)-1 {
            logflow.Printf("main", "Connect %d to %d locally", i, i+1)
            c.ConnectToLocal(self[i+1])
        }
        if i > 0 {
            logflow.Printf("main", "Connect %d to %d locally", i, i-1)
            c.ConnectToLocal(self[i-1])
        }
    }
}

func (self CoordinatorSlice) PrepareAgentProxies() {
    for _, c := range(self) {
        c.RunExporter()
        c.PrepareAgentProxies()
    }
}

func (self CoordinatorSlice) ChainTCP() {
    logflow.Println("main", "Exporting channels")
    for i, c := range(self) {
        if i < len(self)-1 {
            c.ExportRemote(i+1)
        }
        if i > 0 {
            c.ExportRemote(i-1)
        }
    }
    logflow.Println("main", "Connecting coordinators")
    for i, c := range(self) {
        if i < len(self)-1 {
            logflow.Printf("main", "Connect %d to %d over TCP", i, i+1)
            c.ConnectToRPCServer(i+1)
        }
        if i > 0 {
            logflow.Printf("main", "Connect %d to %d over TCP", i, i-1)
            c.ConnectToRPCServer(i-1)
        }
    }
}

func (self CoordinatorSlice) StartAndConnectAgents(agents map[uint32]agent.Agent) {
    for _, coord := range(self) {
        for _, agent_desc := range(coord.conf.Agents) {
            logflow.Print("Starting agent ", agent_desc.Id, " in ", coord.conf)
            go agent.RunWithCoordinator(agents[agent_desc.Id], coord.Address())
        }
    }
}