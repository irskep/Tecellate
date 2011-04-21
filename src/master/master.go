/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: master/master.go
*/

package master

import (
    "agent"
    coordconf "coord/config"
    coordrunner "coord/runner"
    geo "coord/geometry"
    "io/ioutil"
    "json"
    "log"
    "logflow"
    "netchan"
    "time"
    "util"
)

// Config types

type CoordConfig struct {
    Identifier int
    BottomLeft *geo.Point
    TopRight *geo.Point
    Peers []string
    Logs coordconf.LogConfigList
    Agents []*coordconf.AgentDefinition
}

type MasterConfig struct {
    Logs coordconf.LogConfigList
    Coordinators map[string]*CoordConfig
    Agents map[string]*agent.AgentConfig
    MaxTurns int
    MessageStyle string
    UseFood bool
    Size geo.Point
}

// Master

type Master struct {
    conf *MasterConfig
    log logflow.Logger
    coordSendChannels map[string]coordrunner.CoordComm
    coordRecvChannels map[string]coordrunner.CoordComm
    agentSendChannels map[string]agent.AgentComm
    agentRecvChannels map[string]agent.AgentComm
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
    self.conf.fillInData()
    self.log.Print(self.conf.Coordinators)
    self.importCoordChannels()
    self.importAgentChannels()
    self.sendCoordConfigs()
    self.sendAgentConfigs()
    self.sendConnect()
    self.sendGo()
    self.log.Print("Done")
    time.Sleep(1e9/2)
}

func (self *Master) importCoordChannels() {
    self.coordSendChannels = make(map[string]coordrunner.CoordComm, len(self.conf.Coordinators))
    self.coordRecvChannels = make(map[string]coordrunner.CoordComm, len(self.conf.Coordinators))
    for address, _ := range(self.conf.Coordinators) {
        ch_send := make(coordrunner.CoordComm)
        ch_recv := make(coordrunner.CoordComm)

        imp := util.MakeImporterWithRetry("tcp", address, 10, self.log)

        self.log.Print("Importing coord master_req")

    	err := imp.Import("master_req", ch_send, netchan.Send, 1)
    	if err != nil {
    	    self.log.Fatal(err)
    	}

        self.log.Print("Importing coord master_rsp")

    	err = imp.Import("master_rsp", ch_recv, netchan.Recv, 1)
    	if err != nil {
    	    self.log.Fatal(err)
    	}
        
        self.coordSendChannels[address] = ch_send
        self.coordRecvChannels[address] = ch_recv
    }
}

func (self *Master) importAgentChannels() {
    self.agentSendChannels = make(map[string]agent.AgentComm, len(self.conf.Agents))
    self.agentRecvChannels = make(map[string]agent.AgentComm, len(self.conf.Agents))
    for address, _ := range(self.conf.Agents) {
        ch_send := make(agent.AgentComm)
        ch_recv := make(agent.AgentComm)

        imp := util.MakeImporterWithRetry("tcp", address, 10, self.log)

        self.log.Print("Importing agent master_req")

    	err := imp.Import("master_req", ch_send, netchan.Send, 1)
    	if err != nil {
    	    self.log.Fatal(err)
    	}

        self.log.Print("Importing agent master_rsp")

    	err = imp.Import("master_rsp", ch_recv, netchan.Recv, 1)
    	if err != nil {
    	    self.log.Fatal(err)
    	}
        
        self.agentSendChannels[address] = ch_send
        self.agentRecvChannels[address] = ch_recv
    }
}

func (self *Master) sendCoordConfigs() {
    for address, cc := range(self.conf.Coordinators) {
        self.log.Print("Configuring ", address)
        
        peers := make(map[string]int)
        for _, peerAddress := range(cc.Peers) {
            peerConf := self.conf.Coordinators[peerAddress]
            peers[peerAddress] = peerConf.Identifier
        }
        thisConf := coordconf.NewConfig(cc.Identifier,
                                        address,
                                        self.conf.MaxTurns,
                                        cc.Agents,
                                        self.conf.MessageStyle,
                                        self.conf.UseFood,
                                        cc.BottomLeft,
                                        cc.TopRight)
        thisConf.Logs = cc.Logs
        thisConf.Peers = peers
        
        bytes, err := json.Marshal(thisConf)
        if err != nil {
            self.log.Fatal(err)
        }
        self.coordSendChannels[address] <- bytes
        
        self.log.Println("Waiting for response")
        rsp := <- self.coordRecvChannels[address]
        if string(rsp) != "configured" {
            self.log.Fatal("Coordinator at ", address, " failed: ", string(bytes))
        }
    }
}

func (self *Master) sendAgentConfigs() {
    for address, ac := range(self.conf.Agents) {
        self.log.Print("Configuring ", address)
        
        bytes, err := json.Marshal(ac)
        if err != nil {
            self.log.Fatal(err)
        }
        self.agentSendChannels[address] <- bytes
        
        self.log.Println("Waiting for response")
        rsp := <- self.coordRecvChannels[address]   // "ok"
        if string(rsp) != "configured" {
            self.log.Fatal("Coordinator at ", address, " failed: ", string(bytes))
        }
    }
}

func (self *Master) sendConnect() {
    for address, _ := range(self.coordSendChannels) {
        self.log.Print("Sending connect to ", address)
        self.coordSendChannels[address] <- []byte("connect")
    }
    for address, _ := range(self.coordSendChannels) {
        self.log.Print("Waiting for ", address)
        <- self.coordRecvChannels[address]  // "ok"
    }
}

func (self *Master) sendGo() {
    for addr, ch_send := range(self.agentSendChannels) {
        self.log.Print("Sending go to ", addr)
        ch_send <- []byte("go")
    }
    for addr, ch_send := range(self.coordSendChannels) {
        self.log.Print("Sending go to ", addr)
        ch_send <- []byte("go")
    }
}

func (self *MasterConfig) fillInData() {
    currentCoordId := 0
    for coordAddress, coordConf := range(self.Coordinators) {
        currentCoordId += 1
        coordConf.Identifier = currentCoordId
        coordConf.Agents = make([]*coordconf.AgentDefinition, 0)
        bl := coordConf.BottomLeft
        tr := coordConf.TopRight
        var currentAgentId uint32 = 0
        for _, agentConf := range(self.Agents) {
            currentAgentId += 1
            agentConf.Id = currentAgentId
            ap := agentConf.Position
            if bl.X <= ap.X && ap.X < tr.X && bl.Y <= ap.Y && ap.Y < tr.Y {
                ad := coordconf.NewAgentDefinition(agentConf.Id, ap.X, ap.Y, agentConf.Energy)
                agentConf.CoordAddress = coordAddress
                coordConf.Agents = append(coordConf.Agents, ad)
            }
        }
    }
}
