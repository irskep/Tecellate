package wifi

import "testing"

import (
    "agent"
    "agent/link"
    "coord"
    cagent "coord/agent"
    aproxy "coord/agent/proxy"
    geo "coord/geometry"
)

// func AgentFactories(gameconf *coord.GameConfig) {
func makeStaticAgent(id uint, pos *geo.Point, energy cagent.Energy) *aproxy.AgentProxy {
    agnt := make(chan link.Message, 10)
    prox := make(chan link.Message, 10)
    simple := NewStaticBot(id)
    proxy := aproxy.NewAgentProxy(prox, agnt)
    proxy.SetState(cagent.NewAgentState(0, *pos, energy))
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    return proxy
}
// }

func makeRandAgent(id uint, pos *geo.Point, energy cagent.Energy) *aproxy.AgentProxy {
    agnt := make(chan link.Message, 10)
    prox := make(chan link.Message, 10)
    simple := NewRandomBot(id)
    proxy := aproxy.NewAgentProxy(prox, agnt)
    proxy.SetState(cagent.NewAgentState(0, *pos, energy))
    go func() {
        agent.Run(simple, agnt, prox)
    }()
    return proxy
}

func TestStatic(t *testing.T) {
    defer initLogs("TestStatic", t)()

    var time cagent.Energy = 1000
    var gameconf *coord.GameConfig = coord.NewGameConfig(int(time), "noise", true, false, 100, 100)
    gameconf.AddAgent(makeStaticAgent(1, geo.NewPoint(0, 0), time))
    gameconf.AddAgent(makeStaticAgent(2, geo.NewPoint(6, 6), time))
    gameconf.AddAgent(makeStaticAgent(3, geo.NewPoint(12, 12), time))
    gameconf.AddAgent(makeStaticAgent(4, geo.NewPoint(18, 18), time))
    gameconf.AddAgent(makeStaticAgent(5, geo.NewPoint(24, 24), time))
    gameconf.AddAgent(makeStaticAgent(6, geo.NewPoint(30, 30), time))
    gameconf.AddAgent(makeStaticAgent(7, geo.NewPoint(36, 36), time))
    gameconf.AddAgent(makeStaticAgent(8, geo.NewPoint(42, 42), time))
    coords := gameconf.InitWithChainedLocalCoordinators(1, 60)
    coords.Run()
}
