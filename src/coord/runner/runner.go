package runner

import (
    "coord"
    "coord/config"
    "fmt"
    "json"
    "logflow"
    "netchan"
)

type CoordComm chan []byte

type CoordRunner struct {
    myCoord *coord.Coordinator
    masterReq CoordComm
    masterRsp CoordComm
    log logflow.Logger
}

func New(c *coord.Coordinator, address string) *CoordRunner {
    return &CoordRunner{
        myCoord: c,
        masterReq: make(CoordComm),
        masterRsp: make(CoordComm),
        log: logflow.NewSource(fmt.Sprintf("coordrunner/%v", address)),
    }
}

func (self *CoordRunner) ExportNetchans() {
    self.myCoord.Exporter.Export("master_req", self.masterReq, netchan.Recv)
    self.myCoord.Exporter.Export("master_rsp", self.masterReq, netchan.Send)
}

func (self *CoordRunner) RunExporter() {
    self.myCoord.RunExporterInitial()
}

func (self *CoordRunner) ReadConfig() {
    bytes := <- self.masterReq
    configObj := new(config.Config)
    err := json.Unmarshal(bytes, configObj)
    if err != nil {
        self.log.Fatal(err)
    } else {
        self.log.Print("Configured with ", configObj)
    }
    self.myCoord.Configure(configObj)
}
