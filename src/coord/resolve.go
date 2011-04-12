package coord

import "fmt"
import cagent "coord/agent"
// import geo "coord/geometry"

func (self *Coordinator) transformsForNextTurn(peerData []*GameStateResponse) []cagent.Transform {
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))
    
    self.log.Printf("From my neighbors, I see:")
    for _, s := range peerData {
        self.log.Printf("%v", *s)
    }

    exec_turn := func(agent cagent.Agent, done chan<- bool) {
        agent.Turn()
        done <- true
    }

    waiting := make(chan bool, len(agents))
    // for each agent
    //     execute the turn
    for _, agent := range(agents) {
        go exec_turn(agent, waiting)
    }

    // for each agent
    //     ensure the turn has been exec
    for _, _ = range(agents) {
        <-waiting
    }

    // for each agent
    //     construct a StateTransform
    fmt.Println("\n\n---------- Starting Resolve -----------\n")
    moves := make(map[complex128]*StateTransform, len(agents))
    for i, agent := range(agents) {
        state := agent.State()
        t := transformFromState(state)
        t.turn = self.availableGameState.Turn+1

        if state.Wait > 0 {
            t.wait = state.Wait - 1
        } else {
            t.wait = 0
        }

        if state.Inventory().Energy > 0 {
            t.energy = state.Inventory().Energy - 1
            t.alive = true
        } else {
            t.energy = 0
            t.alive = false
        }

        if state.Alive && state.Move != nil {
            t.pos = state.Move.Position.Add(state.Position)
        } else {
            t.pos = state.Position
        }

        if _, has := moves[t.pos.Complex()]; !has {
            moves[t.pos.Complex()] = t
        } else {
            moves[t.pos.Complex()].pos = moves[t.pos.Complex()].state.Position
            t.pos = t.state.Position
        }

        transforms[i] = t
    }

    for _, transform := range(transforms) {
        fmt.Println(transform)
    }


    fmt.Println("\n---------- Ending Resolve -----------\n\n")
    return transforms;
}