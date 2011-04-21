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
    r.ExportMasterConn()
    r.ReadConfigFromMaster()
    r.RunExporter()
    r.FinishConfig()
    r.ConnectToPeers()
    r.WaitForGo()
    r.Run()
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

func (self *CoordRunner) ExportMasterConn() {
    self.myCoord.Exporter.Export("master_req", self.masterReq, netchan.Recv)
    self.myCoord.Exporter.Export("master_rsp", self.masterRsp, netchan.Send)
    self.myCoord.RunExporterBlocking(1)
}

func (self *CoordRunner) ReadConfigFromMaster() {
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

func (self *CoordRunner) RunExporter() {
    self.log.Print(self.myCoord.NumInitialConns())
    go self.myCoord.RunExporterBlocking(self.myCoord.NumInitialConns())
}

func (self *CoordRunner) FinishConfig() {
    self.myCoord.PrepareAgentProxies()
    self.myCoord.PrepareCoordProxies()
    self.masterRsp <- []byte("configured")
}

func (self *CoordRunner) ConnectToPeers() {
    self.log.Print("Ready to connect")
    <- self.masterReq // "connect"
    self.log.Print("Received connect")
    self.myCoord.ConnectCoordProxies()
    self.masterRsp <- []byte("ok")
}

func (self *CoordRunner) WaitForGo() {
    <- self.masterReq
    self.log.Print("Go received")
}

func (self *CoordRunner) Run() {
    self.myCoord.StartRPCServer()
    complete := make(chan bool)
    go self.myCoord.ProcessTurns(complete)
    <- complete
}

func (self *CoordRunner) Close() {
    self.myCoord.Close()
    close(self.masterReq)
    close(self.masterRsp)
}
