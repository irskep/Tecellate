/* Set up various configurations of coordinators */

package coord

import cagent "coord/agent"
import "coord/config"
import geo "coord/geometry"

type GameConfig struct {
    MaxTurns int
    MessageStyle string
    UseFood bool
    Size geo.Point
    Agents []cagent.Agent
}

func NewGameConfig(maxTurns int, msgStyle string, food bool, w, h int) *GameConfig {
    return &GameConfig{MaxTurns: maxTurns,
                       MessageStyle: msgStyle,
                       UseFood: food,
                       Size: *geo.NewPoint(w, h),
                       Agents: make([]cagent.Agent, 0),
    }
}

func (self *GameConfig) AddAgent(a cagent.Agent) {
    self.Agents = append(self.Agents, a)
}

func (self *GameConfig) CoordConfig(id int, bl *geo.Point, tr *geo.Point) *config.Config {
    thisCoordsAgents := make([]cagent.Agent, 0)
    
    for _, a := range self.Agents {
        p := a.State().Position
        if bl.X <= p.X && p.X < tr.X && bl.Y <= p.Y && p.Y < tr.Y {
            thisCoordsAgents = append(thisCoordsAgents, a)
        }
    }
    
    return config.NewConfig(id, 
                            self.MaxTurns,
                            thisCoordsAgents, 
                            self.MessageStyle, 
                            self.UseFood, 
                            bl, 
                            tr)
}

func (self *GameConfig) InitWithSingleLocalCoordinator() *Coordinator {
    c := NewCoordinator()
    c.Configure(self.CoordConfig(0, geo.NewPoint(0, 0), geo.NewPoint(self.Size.X, self.Size.Y)))
    return c
}

func (self *GameConfig) InitWithChainedLocalCoordinators(k int, w int) CoordinatorSlice {
    coords := self.SideBySideCoordinators(k, w, self.Size.Y)
    coords.Chain()
    return coords
}

func (self *GameConfig) InitWithTCPChainedLocalCoordinators(k int, w int) CoordinatorSlice {
    coords := self.SideBySideCoordinators(k, w, self.Size.Y)
    coords.ChainTCP()
    return coords
}

func (self *GameConfig) SideBySideCoordinators(k, w, h int) CoordinatorSlice {
    coords := make(CoordinatorSlice, k)
    for i := 0; i < k; i++ {
        newConf := self.CoordConfig(i, geo.NewPoint(w*i, 0), geo.NewPoint(w*(i+1), h))
        coords[i] = NewCoordinator()
        coords[i].Configure(newConf)
    }
    return coords
}
