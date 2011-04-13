package main

import "logflow"
import "agents/configurable"
// import "agent"

func main() {
    logflow.StdoutSink(".*")
    logflow.Println("main", "agents/simple/main.go")
    cf := configurable.New(0)
    logflow.Println("main", cf)
    // agent.Run(cf)
}
