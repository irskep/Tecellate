package wifi

import "testing"
import "fmt"

import (
    "agent"
    "agent/link"
    "coord"
    cagent "coord/agent"
    aproxy "coord/agent/proxy"
    geo "coord/geometry"
    "agents/wifi/lib"
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

func run_static(time cagent.Energy) (uint32, uint32, []*StaticBot) {
    gameconf := coord.NewGameConfig(int(time), "noise", true, false, 100, 100)
    f := AgentFactories(gameconf)

    var first uint32 = 1
    var last uint32 = 8

    var bots []*StaticBot = []*StaticBot{
        f["Static"](uint(first), geo.NewPoint(0, 0), time).(*StaticBot),
        f["Static"](2, geo.NewPoint(6, 6), time).(*StaticBot),
        f["Static"](3, geo.NewPoint(12, 12), time).(*StaticBot),
        f["Static"](4, geo.NewPoint(18, 18), time).(*StaticBot),
        f["Static"](5, geo.NewPoint(24, 24), time).(*StaticBot),
        f["Static"](6, geo.NewPoint(30, 30), time).(*StaticBot),
        f["Static"](7, geo.NewPoint(36, 36), time).(*StaticBot),
        f["Static"](uint(last), geo.NewPoint(42, 42), time).(*StaticBot),
    }
    coords := gameconf.InitWithChainedLocalCoordinators(1, 60)
    coords.Run()

    return first, last, bots
}

func TestStatic_run10(t *testing.T) {
    _, closer := initLogs("TestStatic_run10", t)
    defer closer()

    run_static(200)
}

func TestStatic_Neighbors(t *testing.T) {
    log, closer := initLogs("TestStatic_Neighbors", t)
    defer closer()

    first, last, bots := run_static(200)

    check := func(id, i uint32, neighbors lib.Neighbors) {
        if !neighbors.In(i) {
            msg := fmt.Sprintf("id %v not in bot %v neighbors %v",
                                i, id, neighbors,
            )
            log(msg)
            t.Fatal(msg)
        }
    }

    for _, bot := range bots {
        id := uint32(bot.Id())
        neighbors := bot.hello.Neighbors()
        if id != first { check(id, id - 1, neighbors) }
        if id != last { check(id, id + 1, neighbors)}
    }
}

