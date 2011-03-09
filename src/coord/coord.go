/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/coord.go
*/

package coord

type Coordinator struct {
    AvailableGameState GameState
    Peers []CoordinatorProxy
}

type CoordinatorProxy struct {
    
}
