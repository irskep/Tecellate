/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/main.go

COORD MAIN
*/

package main

import (
    "fmt"
    "coord"
    "log"
)

func main() {
    fmt.Println("coord/main.go")
    
    a := coord.NewCoordinator()
    b := coord.NewCoordinator()
    
    a.ConnectToLocal(b)
    b.ConnectToLocal(a)
    
    a.StartRPCServer()
    b.StartRPCServer()
    
    complete := make(chan bool)
    
    go a.ProcessTurns(complete)
    go b.ProcessTurns(complete)
    
    <- complete
    <- complete
    
    log.Println("Done")
}
