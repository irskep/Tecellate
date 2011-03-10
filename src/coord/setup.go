/* Set up various configurations of coordinators */

package coord

import (
    "coord/config"
    "log"
)

func CoordinatorList(k int, configTemplate *config.Config) []*Coordinator {
    coords := make([]*Coordinator, k)
    for i := 0; i < k; i++ {
        coords[i] = NewCoordinator()
        coords[i].Configure(&config.Config{i,
                                           configTemplate.AgentStarts, 
                                           configTemplate.MessageStyle, 
                                           configTemplate.UseFood, 
                                           configTemplate.RandomlyDelayProcessing})
    }
    return coords
}

func ConnectInChain(coords []*Coordinator) {
    for i, c := range(coords) {
        if i < len(coords)-1 {
            log.Printf("main: Connect %d to %d", i, i+1)
            c.ConnectToLocal(coords[i+1])
        }
        if i > 0 {
            log.Printf("main: Connect %d to %d", i, i-1)
            c.ConnectToLocal(coords[i-1])
        }
    }
}
