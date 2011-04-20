/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: master/master.go
*/

package master

import (
    coordconf "coord/config"
    "io/ioutil"
    "json"
    "log"
    "logflow"
    "netchan"
    "util"
)

// Config types

type LogConfig []string
type LogConfigList []LogConfig

type CoordConfig struct {
    BottomLeft []int
    TopRight []int
    Peers []string
    Logs LogConfigList
    Agents []*coordconf.AgentDefinition
}

type AgentConfig struct {
    Id uint32
    Position []int
    Energy int
    Logs LogConfigList
}

type MasterConfig struct {
    Logs LogConfigList
    Coordinators map[string]CoordConfig
    Agents map[string]AgentConfig
}

// Master

type CoordComm chan []byte

type Master struct {
    conf *MasterConfig
    log logflow.Logger
    coordSendChannels map[string]CoordComm
    coordRecvChannels map[string]CoordComm
}

func New(args []string) *Master {
    mc := new(MasterConfig)
    
    txt, err := ioutil.ReadFile(args[1])
    if err != nil {
        log.Fatal(err)
    }
    err = json.Unmarshal(txt, mc)
    if err != nil {
        log.Fatal(err)
    }
    
    m := &Master{
        conf: mc,
        log: logflow.NewSource("master"),
    }
    
    m.conf.Logs.Apply()
    
    m.log.Print("Configured.")
    
    return m
}

func (self *Master) ConnectToCoords() {
    self.conf.fillInAgentLists()
    self.importCoordChannels()
    self.sendCoordConfigs()
    self.ackCoordConfigs()
    self.sendGo()
}

func (self *Master) importCoordChannels() {
    self.coordSendChannels = make(map[string]CoordComm, len(self.conf.Coordinators))
    self.coordRecvChannels = make(map[string]CoordComm, len(self.conf.Coordinators))
    for address, _ := range(self.conf.Coordinators) {
        ch_send := make(CoordComm)
        ch_recv := make(CoordComm)

        imp := util.MakeImporterWithRetry("tcp", address, 10, self.log)

        self.log.Print("Importing master_req")

    	err := imp.Import("master_req", ch_send, netchan.Send, 1)
    	if err != nil {
    	    self.log.Fatal(err)
    	}

        self.log.Print("Importing master_rsp")

    	err = imp.Import("master_rsp", ch_recv, netchan.Recv, 1)
    	if err != nil {
    	    self.log.Fatal(err)
    	}
        
        self.coordSendChannels[address] = ch_send
        self.coordRecvChannels[address] = ch_recv
    }
}

func (self *Master) sendCoordConfigs() {
    for address, ch_send := range(self.coordSendChannels) {
        bytes, err := json.Marshal(self.conf.Coordinators[address])
        if err != nil {
            self.log.Fatal(err)
        }
        ch_send <- bytes
    }
}

func (self *Master) ackCoordConfigs() {
    for address, ch_recv := range(self.coordSendChannels) {
        bytes := <- ch_recv
        if string(bytes) != "configured" {
            self.log.Fatal("Coordinator at ", address, " failed: ", string(bytes))
        }
    }
}

func (self *Master) sendGo() {
    for _, ch_send := range(self.coordSendChannels) {
        ch_send <- []byte("go")
    }
}

func (self *MasterConfig) fillInAgentLists() {
    for _, coordConf := range(self.Coordinators) {
        coordConf.Agents = make([]*coordconf.AgentDefinition, 0)
        blx := coordConf.BottomLeft[0]
        bly := coordConf.BottomLeft[1]
        trx := coordConf.TopRight[0]
        try := coordConf.TopRight[1]
        var currentId uint32 = 0
        for _, agentConf := range(self.Agents) {
            currentId += 1
            agentConf.Id = currentId
            ax := agentConf.Position[0]
            ay := agentConf.Position[1]
            if blx <= ax && ax < trx && bly <= ay && ay < try {
                ad := coordconf.NewAgentDefinition(agentConf.Id, ax, ay, agentConf.Energy)
                coordConf.Agents = append(coordConf.Agents, ad)
            }
        }
    }
}

// Logs

func (self LogConfigList) Apply() {
    for _, l := range(self) {
        l.Apply()
    }
}

func (self LogConfig) Apply() {
    switch self[0] {
    case "stdout":
        logflow.StdoutSink(self[1])
    case "file":
        logflow.FileSink(self[1], true, self[2:]...)
    }
}
