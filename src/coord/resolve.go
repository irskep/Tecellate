package coord

import cagent "coord/agent"
import aproxy "coord/agent/proxy"
import game "coord/game"
// import geo "coord/geometry"

func (self *Coordinator) getNewAgents(peers []*game.GameStateResponse) []cagent.Agent {
    newAgents := make([]cagent.Agent, 0)
    for _, rsp := range(peers) {
        if len(rsp.AgentsToAdopt) > 0 {
            self.log.Print("Got new agents from ", rsp.Identifier)
            for _, as := range(rsp.AgentsToAdopt) {
                self.log.Print("One of them is ", as)
                // newAgents = append(newAgents, self.NewProxy(as))
                self.AddNewProxyFromState(&as)
            }
            self.RunExporterBlocking(len(rsp.AgentsToAdopt))
        }
    }
    return newAgents
}

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

func (self *Coordinator) transformsForNextTurn(peers []*game.GameStateResponse) ([]cagent.Transform, *game.Messages, *game.Messages, []cagent.Agent) {
    messages := game.NewMessages(peers)
    myMessages := game.NewMessages(nil)
    
    newAgents := self.getNewAgents(peers)
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))
    
    self.doTurns(agents)

    // for each agent
    //     construct a StateTransform
    self.log.Println("\n\n---------- Starting Resolve -----------\n")
    
    moves := make(map[complex128]uint32, len(agents))
    for _, peerGameState := range(peers) {
        for _, st := range peerGameState.AgentStates {
            moves[st.Position.Complex()] = st.Id
            if st.Move.Valid {
                requestedPosition := st.Move.Position.Add(st.Position)
                moves[requestedPosition.Complex()] = st.Id
            }
        }
    }
    
    for _, agent := range(agents) {
        moves[agent.State().Position.Complex()] = agent.State().Id
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
        
        if self.conf.UseFood {
            if state.Energy > 0 {
                t.energy = state.Energy - 1
                t.alive = true
            } else {
                t.energy = 0
                t.alive = false
            }
        }

        if state.Alive && state.Move.Valid {
            requestedPosition := *state.Move.Position.Add(state.Position)
            if occupant, has := moves[requestedPosition.Complex()]; occupant != state.Id && has {
                self.log.Print("Agent ", state.Id, " fails move ", state.Position, " - ", requestedPosition)
                t.pos = state.Position
            } else {
                self.log.Print("Agent ", state.Id, " performs move ", state.Position, " - ", requestedPosition)
                moves[requestedPosition.Complex()] = state.Id
                t.pos = requestedPosition
            }
            
            for _, msg := range state.Move.Messages {
                messages.Add(msg)
                myMessages.Add(msg)
            }
        }

        transforms[i] = t
    }

    self.log.Println("\n---------- Ending Resolve -----------\n\n")
    return transforms, messages, myMessages, newAgents
}
