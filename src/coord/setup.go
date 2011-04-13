/* Set up various configurations of coordinators */

package coord

import "coord/config"
import geo "coord/geometry"

import (
    "logflow"
)

/* Stick together building blocks */

func ChainedLocalCoordinators(k int, configTemplate *config.Config, w int, h int) CoordinatorSlice {
    coords := CoordinatorList(k, configTemplate, geo.NewPoint(w, h))
    ConnectInChain(coords)
    return coords
}

/* Building blocks for making coordinators */

func CoordinatorList(k int, configTemplate *config.Config, size *geo.Point) CoordinatorSlice {
    coords := make(CoordinatorSlice, k)
    w := size.X/k
    h := size.Y
    for i := 0; i < k; i++ {
        newConf := configTemplate.Duplicate(i, geo.NewPoint(w*i, 0), geo.NewPoint(w*(i+1), h))
        coords[i] = NewCoordinator()
        coords[i].Configure(newConf)
    }
    return coords
}

/* Building blocks for connecting coordinators */

func ConnectInChain(coords CoordinatorSlice) {
    for i, c := range(coords) {
        if i < len(coords)-1 {
            logflow.Printf("main", "Connect %d to %d", i, i+1)
            c.ConnectToLocal(coords[i+1])
        }
        if i > 0 {
            logflow.Printf("main", "Connect %d to %d", i, i-1)
            c.ConnectToLocal(coords[i-1])
        }
    }
}
