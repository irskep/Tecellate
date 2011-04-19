package wifi

// import "testing"
// import "fmt"

import "agent"
// import "agent/link"
import "coord"
// import cagent "coord/agent"
// import aproxy "coord/agent/proxy"
// import geo "coord/geometry"

type AgentFactory func(uint32, int, int, int) agent.Agent

func AgentFactories(gameconf *coord.GameConfig, first, last uint32) map[string]AgentFactory {
    return map[string]AgentFactory {
    "Static":
        func(id uint32, x, y, energy int) agent.Agent {
            bot := NewStaticBot(id, first, last)
            gameconf.AddAgent(id, x, y, energy)
            return bot
        },
//     "Random":
//         func(id uint, pos *geo.Point, energy cagent.Energy) agent.Agent {
//             agnt := make(chan link.Message, 10)
//             prox := make(chan link.Message, 10)
//             bot := NewRandomBot(id)
//             proxy := aproxy.NewAgentProxy(prox, agnt)
//             proxy.SetState(cagent.NewAgentState(0, *pos, energy))
//             go func() {
//                 agent.Run(bot, agnt, prox)
//             }()
//             gameconf.AddAgent(proxy)
//             return bot
//         },
    }
}

func run_static(time int) (uint32, uint32, []*StaticBot) {
    gameconf := coord.NewGameConfig(int(time), "noise", true, 100, 100)
    var first uint32 = 1
    var last uint32 = 8

    f := AgentFactories(gameconf, first, last)

    var bots []*StaticBot = []*StaticBot{
        f["Static"](first, 0, 0, time).(*StaticBot),
        f["Static"](2, 6, 6, time).(*StaticBot),
        f["Static"](3, 12, 12, time).(*StaticBot),
        f["Static"](4, 18, 18, time).(*StaticBot),
        f["Static"](5, 24, 24, time).(*StaticBot),
        f["Static"](6, 30, 30, time).(*StaticBot),
        f["Static"](7, 36, 36, time).(*StaticBot),
        f["Static"](last, 42, 42, time).(*StaticBot),
    }
    agents := make(map[uint32]agent.Agent)
    for _, bot := range bots {
        agents[bot.Id()] = agent.Agent(bot)
    }
    coords := gameconf.InitWithChainedLocalCoordinators(1, agents)
    coords.Run()

    return first, last, bots
}
