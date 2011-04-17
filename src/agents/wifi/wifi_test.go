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
    "agents/wifi/testerlib"
)


type Neighbors []uint32

func (self Neighbors) In(id uint32) bool {
    for _, cur := range self {
        if id == cur {
            return true
        }
    }
    return false
}

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

func TestStatic_run8(t *testing.T) {
    _, closer := testerlib.InitLogs("TestStatic_run8", t)
    defer closer()

    run_static(200)
}

func check(id, i uint32, neighbors Neighbors) (string, bool) {
    if !neighbors.In(i) {
        msg := fmt.Sprintf("id %v not in bot %v neighbors %v", i, id, neighbors)
        return msg, false
    }
    return "", true
}

func TestStatic_Neighbors(t *testing.T) {
    _, closer := testerlib.InitLogs("TestStatic_Neighbors", t)
    defer closer()

    first, last, bots := run_static(200)

    for _, bot := range bots {
        id := uint32(bot.Id())
        neighbors := bot.hello.Neighbors()
        if id != first {
            if msg, ok := check(id, id - 1, neighbors); !ok {
                t.Error(msg)
            }
        }
        if id != last {
            if msg, ok := check(id, id + 1, neighbors); !ok {
                t.Error(msg)
            }
        }
    }
}

func TestStatic_Reachable(t *testing.T) {
    log, closer := testerlib.InitLogs("TestStatic_Reachable", t)
    defer closer()

    var msgs []string
    var success bool
    for j := 0; j < 5; j++ {
        first, last, bots := run_static(500)
        msgs = make([]string, 0, len(bots)*len(bots))
        success = true
        for _, bot := range bots {
            id := uint32(bot.Id())
            reachable := bot.route.Reachable()
            for i := first; i <= last; i++ {
                if msg, ok := check(id, i, reachable); !ok {
                    msgs = append(msgs, msg)
                    success = false
                }
            }
        }
        if success { break }
        log("the routing tables where not complete, about to retry")
        log("errors")
        for _, msg := range msgs {
            log("    ", msg)
        }
        log("retrying...\n")
    }
    if !success {
        t.Error("The routing tables where never completed...")
    }
}

