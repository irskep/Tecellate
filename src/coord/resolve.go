package coord

import cagent "coord/agent"
import aproxy "coord/agent/proxy"
import game "coord/game"
// import geo "coord/geometry"

func (self *Coordinator) doTurns(agents []cagent.Agent) {
    exec_turn := func(agent cagent.Agent, done chan<- bool) {
        agent.(*aproxy.AgentProxy).SetGameState(self.availableGameState)
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
}

func (self *Coordinator) transformsForNextTurn(peers []*game.GameStateResponse) ([]cagent.Transform, *game.Messages) {
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))
    messages := game.NewMessages(peers)

    self.log.Println("From my neighbors, I see:", peers)
    for _, s := range peers {
        self.log.Printf("    %vln", *s)
    }
    
    self.doTurns(agents)

    // ---------------------------------------------------------------------
    //TODO:
    //  Iterate over peer data to resolve peer turns.
    // ---------------------------------------------------------------------

    // for each agent
    //     construct a StateTransform
    self.log.Println("\n\n---------- Starting Resolve -----------\n")
    
    moves := make(map[complex128]int, len(agents))
    for _, peerGameState := range(peers) {
        for _, st := range peerGameState.AgentStates {
            moves[st.Position.Complex()] = 1
            moves[st.Move.Position.Complex()] = 1
        }
    }
    
    for _, agent := range(agents) {
        moves[agent.State().Position.Complex()] = 1
    }
    
    for i, agent := range(agents) {
        state := agent.State()
        t := transformFromState(state)
        t.turn = self.availableGameState.Turn+1

        if state.Wait > 0 {
            t.wait = state.Wait - 1
        } else {
            t.wait = 0
        }

        if state.Energy > 0 {
            t.energy = state.Energy - 1
            t.alive = true
        } else {
            t.energy = 0
            t.alive = false
        }

        if state.Alive && state.Move != nil {
            t.pos = *state.Move.Position.Add(state.Position)
            for _, msg := range state.Move.Messages {
                messages.Add(msg)
            }
        } else {
            t.pos = state.Position
        }
        
        if _, has := moves[t.pos.Complex()]; has {
            t.pos = state.Position
        } else {
            moves[t.pos.Complex()] = 1
        }

        transforms[i] = t
    }

    for _, transform := range(transforms) {
        self.log.Println(transform)
    }
    self.log.Println(messages)


    self.log.Println("\n---------- Ending Resolve -----------\n\n")
    return transforms, messages;
}
