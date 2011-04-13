package configurable

import "agent"

type Configurable struct {
    id uint
    XVelocity int
    YVelocity int
    LogBroadcast bool
    LogMove bool
    LogListen bool
    LogCollect bool
    LogLook bool
    LogPrevious bool
    LogEnergy bool
}

func New(id uint) *Configurable {
    return &Configurable{id, 0, 0, false, false, true, false, false, false, false}
}

func (self *Configurable) Turn(comm agent.Comm) {
    broadcasted := comm.Broadcast(23, []byte("hello_world"))
    if (self.LogBroadcast) {
        comm.Log("Broadcast success:", broadcasted)
    }
    if !comm.Move(self.XVelocity, self.YVelocity) && self.LogMove {
        comm.Log("Move failed")
    }
    if (self.LogListen) {
        comm.Log("Heard:", string(comm.Listen(23)))
    }
    if (self.LogCollect) {
        comm.Log("Collected: ", comm.Collect())
    }
    if (self.LogLook) {
        comm.Log("Look: ", comm.Collect())
    }
    if (self.LogPrevious) {
        comm.Log("Previous: ", comm.PrevResult())
    }
    if (self.LogEnergy) {
        comm.Log("Energy:", comm.Energy())
    }    
    return
}

func (self *Configurable) Id() uint {
    return self.id
}
