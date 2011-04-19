package wifi

import "testing"


import "fmt"

import (
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

func TestStatic_run8(t *testing.T) {
    defer testerlib.InitLogs("TestStatic_run8", t)()

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
    defer testerlib.InitLogs("TestStatic_Neighbors", t)()

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
//     t.Fatal("lala")
}

func TestStatic_Reachable(t *testing.T) {
    defer testerlib.InitLogs("TestStatic_Reachable", t)()

    var msgs []string
    var success bool
    var confirm_rates []float64
    var misses uint
    for j := 0; j < 5; j++ {
        first, last, bots := run_static(750)
        msgs = make([]string, 0, len(bots)*len(bots))
        success = true
        var confirm_rate float64
        for _, bot := range bots {
            id := uint32(bot.Id())
            reachable := bot.route.Reachable()
            for i := first; i <= last; i++ {
                if msg, ok := check(id, i, reachable); !ok {
                    misses += 1
                    msgs = append(msgs, msg)
                    success = false
                }
            }
            confirm_rate += bot.route.ConfirmRate()
        }
        confirm_rates = append(confirm_rates, confirm_rate/float64(len(bots)))
        if success { break }
        t.Log("the routing tables where not complete, about to retry")
        t.Log("errors")
        for _, msg := range msgs {
            t.Log("    ", msg)
        }
        t.Log("retrying...\n")
    }
    if !success {
        t.Error("The routing tables where never completed...")
        var acc float64
        for _, rate := range confirm_rates {
            acc += rate
            t.Error("Confirm Rate =", rate)
        }
        t.Error("Avg Rate =", acc/float64(len(confirm_rates)))
        t.Error("Misses =", misses)
    }
}

