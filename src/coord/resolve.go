package coord

// import geo "coord/geometry"
import cagent "coord/agent"

func (self *Coordinator) transformsForNextTurn(peerData []*GameStateResponse) []cagent.Transform {
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))

    for ix, agent := range(agents) {
        _ = agent.Turn()
        state := agent.State()
        t := transformFromState(state)
        t.turn = self.availableGameState.Turn+1
        transforms[ix] = t
    }
    return transforms;
}
