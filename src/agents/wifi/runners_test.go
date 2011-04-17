package wifi

// import "testing"
// import "fmt"

import "agent"
import "agent/link"
import "coord"
import cagent "coord/agent"
import aproxy "coord/agent/proxy"
import geo "coord/geometry"

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
