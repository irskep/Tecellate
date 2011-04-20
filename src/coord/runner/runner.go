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

func RunAtAddress(address string) {
    c := coord.NewCoordinator()
    r := New(c, address)
    r.ExportNetchans()
    r.RunExporter()
    r.ReadConfig()
    r.WaitForGo()
    r.log.Print("Done")
    r.Close()
}

func New(c *coord.Coordinator, address string) *CoordRunner {
    r := &CoordRunner{
        myCoord: c,
        masterReq: make(CoordComm),
        masterRsp: make(CoordComm),
        log: logflow.NewSource(fmt.Sprintf("coordrunner/%v", address)),
    }
    c.Config().Address = address
    return r
}

func (self *CoordRunner) ExportNetchans() {
    self.myCoord.Exporter.Export("master_req", self.masterReq, netchan.Recv)
    self.myCoord.Exporter.Export("master_rsp", self.masterRsp, netchan.Send)
}

func (self *CoordRunner) RunExporter() {
    self.myCoord.RunExporterBlocking(self.myCoord.NumInitialConns()+1)
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
    self.masterRsp <- []byte("configured")
}

func (self *CoordRunner) WaitForGo() {
    <- self.masterReq
    self.log.Print("Go received")
}

func (self *CoordRunner) Close() {
    self.myCoord.Close()
    close(self.masterReq)
    close(self.masterRsp)
}
