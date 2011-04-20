/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: master/master.go
*/

package master

import (
    "io/ioutil"
    "json"
    "log"
    "logflow"
)

// Config types

type LogConfig []string
type LogConfigList []LogConfig

type ProcessConfigMap map[string]LogConfigList

type MasterConfig struct {
    Address string
    Logs LogConfigList
    Coordinators ProcessConfigMap
    Agents ProcessConfigMap
}

// Master

type Master struct {
    conf *MasterConfig
    log logflow.Logger
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
