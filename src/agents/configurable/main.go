package main

import (
    "agent"
    "agents/configurable"
    "logflow"
    "os"
)

func main() {
    logflow.StdoutSink(".*")
    a := configurable.New(0)
    a.XVelocity = 1
    agent.RunStandalone(os.Args[1], a)
}
