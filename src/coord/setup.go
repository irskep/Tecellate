/* Set up various configurations of coordinators */

package coord

import (
    "coord/config"
    "log"
)

/* Stick together building blocks */

func ChainedLocalCoordinators(k int, configTemplate *config.Config) CoordinatorSlice {
    coords := CoordinatorList(3, &config.Config{0, nil, "none", false, true})
    ConnectInChain(coords)
    return coords
}

/* Building blocks for making coordinators */

func CoordinatorList(k int, configTemplate *config.Config) CoordinatorSlice {
    coords := make(CoordinatorSlice, k)
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

/* Building blocks for connecting coordinators */

func ConnectInChain(coords CoordinatorSlice) {
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
