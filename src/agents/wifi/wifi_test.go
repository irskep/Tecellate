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

type AgentFactory func(uint, *geo.Point, cagent.Energy) agent.Agent

func AgentFactories(gameconf *coord.GameConfig) map[string]AgentFactory {
    return map[string]AgentFactory {
    "Static":
        func(id uint, pos *geo.Point, energy cagent.Energy) agent.Agent {
            agnt := make(chan link.Message, 10)
            prox := make(chan link.Message, 10)
            bot := NewStaticBot(id)
            proxy := aproxy.NewAgentProxy(prox, agnt)
            proxy.SetState(cagent.NewAgentState(0, *pos, energy))
            go func() {
                agent.Run(bot, agnt, prox)
            }()
            gameconf.AddAgent(proxy)
            return bot
        },
    "Random":
        func(id uint, pos *geo.Point, energy cagent.Energy) agent.Agent {
            agnt := make(chan link.Message, 10)
            prox := make(chan link.Message, 10)
            bot := NewRandomBot(id)
            proxy := aproxy.NewAgentProxy(prox, agnt)
            proxy.SetState(cagent.NewAgentState(0, *pos, energy))
            go func() {
                agent.Run(bot, agnt, prox)
            }()
            gameconf.AddAgent(proxy)
            return bot
        },
    }
}

func TestStatic_run10(t *testing.T) {
    _, closer := initLogs("TestStatic_run10", t)
    defer closer()

    var time cagent.Energy = 1000
    gameconf := coord.NewGameConfig(int(time), "noise", true, false, 100, 100)
    f := AgentFactories(gameconf)
    f["Static"](1, geo.NewPoint(0, 0), time)
    f["Static"](2, geo.NewPoint(6, 6), time)
    f["Static"](3, geo.NewPoint(12, 12), time)
    f["Static"](4, geo.NewPoint(18, 18), time)
    f["Static"](5, geo.NewPoint(24, 24), time)
    f["Static"](6, geo.NewPoint(30, 30), time)
    f["Static"](7, geo.NewPoint(36, 36), time)
    f["Static"](8, geo.NewPoint(42, 42), time)
    coords := gameconf.InitWithChainedLocalCoordinators(1, 60)
    coords.Run()
}

func TestStatic_Neighbors(t *testing.T) {
    log, closer := initLogs("TestStatic_Neighbors", t)
    defer closer()

    var time cagent.Energy = 1000
    gameconf := coord.NewGameConfig(int(time), "noise", true, false, 100, 100)
    f := AgentFactories(gameconf)
    b1 := f["Static"](1, geo.NewPoint(0, 0), time).(*StaticBot)
    f["Static"](2, geo.NewPoint(6, 6), time)
    f["Static"](3, geo.NewPoint(12, 12), time)
    f["Static"](4, geo.NewPoint(18, 18), time)
    f["Static"](5, geo.NewPoint(24, 24), time)
    f["Static"](6, geo.NewPoint(30, 30), time)
    f["Static"](7, geo.NewPoint(36, 36), time)
    f["Static"](8, geo.NewPoint(42, 42), time)
    coords := gameconf.InitWithChainedLocalCoordinators(1, 60)
    coords.Run()
    log(b1)
}

