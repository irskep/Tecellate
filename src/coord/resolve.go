package coord

import "fmt"
// import geo "coord/geometry"
import cagent "coord/agent"

func (self *Coordinator) transformsForNextTurn(peerData []*GameStateResponse) []cagent.Transform {
    agents := self.availableGameState.Agents
    transforms := make([]cagent.Transform, len(agents))
    
    self.log.Printf("From my neighbors, I see:")
    for _, s := range peerData {
        self.log.Printf("%v", *s)
    }

    for _, agent := range(agents) {
        if ok := agent.Turn(); !ok { continue; }
    }
    fmt.Println("\n\n---------- Starting Resolve -----------\n")
    for ix, agent := range(agents) {
        state := agent.State()
        t := transformFromState(state)
        t.turn = self.availableGameState.Turn+1
        if state.Move != nil {
            self.log.Println(state.Move)
        }
        transforms[ix] = t
    }
    fmt.Println("\n---------- Ending Resolve -----------\n\n")
    return transforms;
}
