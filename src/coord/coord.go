/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

import "coord.game"

type Coordinator struct {
    AvailableGameState game.GameState
    Peers []CoordinatorProxy
}

type CoordinatorProxy struct {
    
}
