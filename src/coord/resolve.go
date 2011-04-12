package coord

import "fmt"
import cagent "coord/agent"
// import geo "coord/geometry"

func (self *Coordinator) transformsForNextTurn(peerData []*GameStateResponse) []cagent.Transform {
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))

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
        if state.Move != nil {
            self.log.Println(state.Move)
            t.pos = state.Move.Position.Add(state.Position)
            if _, has := moves[t.pos.Complex()]; !has {
                moves[t.pos.Complex()] = t
            } else {
                moves[t.pos.Complex()].pos = nil
                t.pos = nil
            }
        }
        transforms[i] = t
    }

    // validate the tranforms are non-conflicting

    fmt.Println("\n---------- Ending Resolve -----------\n\n")
    return transforms;
}
