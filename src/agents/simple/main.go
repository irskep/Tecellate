/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: agents/agent1.go
*/

package main

import "fmt"
import "agents/simple"
// import "agent"

func main() {
    fmt.Println("agents/simple/main.go")
    fmt.Println(simple.NewSimple(0))
//     agent.Run(simple.NewSimple())
}
