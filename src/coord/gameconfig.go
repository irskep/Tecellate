/* Set up various configurations of coordinators */

package coord

import "agent"
import "coord/config"
import geo "coord/geometry"

type GameConfig struct {
    MaxTurns int
    MessageStyle string
    UseFood bool
    Size geo.Point
    Agents []*config.AgentDefinition
}

func NewGameConfig(maxTurns int, msgStyle string, food bool, w, h int) *GameConfig {
    return &GameConfig{MaxTurns: maxTurns,
                       MessageStyle: msgStyle,
                       UseFood: food,
                       Size: *geo.NewPoint(w, h),
                       Agents: make([]*config.AgentDefinition, 0),
    }
}

func (self *GameConfig) AddAgent(id uint, x, y int) {
    self.Agents = append(self.Agents, config.NewAgentDefinition(id, x, y))
}

func (self *GameConfig) CoordConfig(id int, bl *geo.Point, tr *geo.Point) *config.Config {
    thisCoordsAgents := make([]*config.AgentDefinition, 0)
    
    for _, ad := range self.Agents {
        if bl.X <= ad.X && ad.X < tr.X && bl.Y <= ad.Y && ad.Y < tr.Y {
            thisCoordsAgents = append(thisCoordsAgents, ad)
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

func (self *GameConfig) InitWithChainedLocalCoordinators(k int, agents map[uint]agent.Agent) CoordinatorSlice {
    coords := self.SideBySideCoordinators(k, self.Size.X/k, self.Size.Y)
    coords.Chain()
    coords.ConnectToLocalAgents(agents)
    return coords
}

func (self *GameConfig) InitWithTCPChainedLocalCoordinators(k int, agents []agent.Agent) CoordinatorSlice {
    coords := self.SideBySideCoordinators(k, self.Size.X/k, self.Size.Y)
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
