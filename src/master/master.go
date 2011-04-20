/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: master/master.go
*/

package master

import (
    coordconf "coord/config"
    geo "coord/geometry"
    "io/ioutil"
    "json"
    "log"
    "logflow"
    "netchan"
    "util"
)

// Config types

type CoordConfig struct {
    BottomLeft *geo.Point
    TopRight *geo.Point
    Peers []*string
    Logs coordconf.LogConfigList
    Agents []*coordconf.AgentDefinition
}

type AgentConfig struct {
    Id uint32
    Position geo.Point
    Energy int
    Logs coordconf.LogConfigList
}

type MasterConfig struct {
    Logs coordconf.LogConfigList
    Coordinators map[string]CoordConfig
    Agents map[string]AgentConfig
    MaxTurns int
    MessageStyle string
    UseFood bool
    Size geo.Point
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
    var currentId int = 0
    for address, ch_send := range(self.coordSendChannels) {
        currentId += 1
        
        cc := self.conf.Coordinators[address]
        thisConf := coordconf.NewConfig(currentId,
                                        address,
                                        self.conf.MaxTurns,
                                        cc.Agents,
                                        self.conf.MessageStyle,
                                        self.conf.UseFood,
                                        cc.BottomLeft,
                                        cc.TopRight)
        
        bytes, err := json.Marshal(thisConf)
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
        bl := coordConf.BottomLeft
        tr := coordConf.TopRight
        var currentId uint32 = 0
        for _, agentConf := range(self.Agents) {
            currentId += 1
            agentConf.Id = currentId
            ap := agentConf.Position
            if bl.X <= ap.X && ap.X < tr.X && bl.Y <= ap.Y && ap.Y < tr.Y {
                ad := coordconf.NewAgentDefinition(agentConf.Id, ap.X, ap.Y, agentConf.Energy)
                coordConf.Agents = append(coordConf.Agents, ad)
            }
        }
    }
}
