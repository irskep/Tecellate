package configurable

import "agent"

type Configurable struct {
    id uint
    xVelocity int
    yVelocity int
    logBroadcast bool
    logMove bool
    logListen bool
    logCollect bool
    logLook bool
    logPrevious bool
    logEnergy bool
}

func New(id uint) *Configurable {
    return &Configurable{id, 0, 0, false, false, true, false, false, false, false}
}

func (self *Configurable) Turn(comm agent.Comm) {
    broadcasted := comm.Broadcast(23, []byte("hello_world"))
    if (self.logBroadcast) {
        comm.Log("Broadcast success:", broadcasted)
    }
    if !comm.Move(self.xVelocity, self.yVelocity) && self.logMove {
        comm.Log("Move failed")
    }
    if (self.logListen) {
        comm.Log("Heard:", string(comm.Listen(23)))
    }
    if (self.logCollect) {
        comm.Log("Collected: ", comm.Collect())
    }
    if (self.logLook) {
        comm.Log("Look: ", comm.Collect())
    }
    if (self.logPrevious) {
        comm.Log("Previous: ", comm.PrevResult())
    }
    if (self.logEnergy) {
        comm.Log("Energy:", comm.Energy())
    }    
    return
}

func (self *Configurable) Id() uint {
    return self.id
}
