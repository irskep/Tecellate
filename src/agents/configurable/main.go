package main

import (
    "agent"
    "agents/configurable"
    "logflow"
    "os"
)

func main() {
    logflow.StdoutSink(".*")
    agent.RunStandalone(os.Args[1], configurable.New(0))
}
