/* Set up various configurations of coordinators */

package coord

import "coord/agent"
import "coord/config"
import geo "coord/geometry"

type GameConfig struct {
    MessageStyle string
    UseFood bool
    RandomlyDelayProcessing bool
    Size geo.Point
    Agents []agent.Agent
}

func NewGameConfig(msgStyle string, food bool, delay bool, w, h int) *GameConfig {
    return &GameConfig{MessageStyle: msgStyle,
                       UseFood: food,
                       RandomlyDelayProcessing: delay,
                       Size: *geo.NewPoint(w, h),
                       Agents: make([]agent.Agent, 0),
    }
}

func (self *GameConfig) AddAgent(a agent.Agent) {
    self.Agents = append(self.Agents, a)
}

func (self *GameConfig) CoordConfig(id int, bl *geo.Point, tr *geo.Point) *config.Config {
    thisCoordsAgents := make([]agent.Agent, 0)
    
    for _, a := range self.Agents {
        p := a.State().Position
        if bl.X <= p.X && p.X < tr.X && bl.Y <= p.Y && p.Y < tr.Y {
            thisCoordsAgents = append(thisCoordsAgents, a)
        }
    }
    
    return config.NewConfig(id, 
                            thisCoordsAgents, 
                            self.MessageStyle, 
                            self.UseFood, 
                            self.RandomlyDelayProcessing, 
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

func (self *GameConfig) SideBySideCoordinators(k, w, h int) CoordinatorSlice {
    coords := make(CoordinatorSlice, k)
    for i := 0; i < k; i++ {
        newConf := self.CoordConfig(i, geo.NewPoint(w*i, 0), geo.NewPoint(w*(i+1), h))
        coords[i] = NewCoordinator()
        coords[i].Configure(newConf)
    }
    return coords
}
