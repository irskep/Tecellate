/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/main.go

COORD MAIN
*/

package main

import (
    "coord/config"
    "fmt"
    "coord"
    "log"
)

func main() {
    fmt.Println("coord/main.go")
    
    // Initialize
    a := coord.NewCoordinator()
    b := coord.NewCoordinator()
    
    conf := config.Config{nil, "boolean", false, true}
    a.Configure(&conf)
    b.Configure(&conf)
    
    // Set up test environment
    a.ConnectToLocal(b)
    b.ConnectToLocal(a)
    
    // Start RPC threads
    a.StartRPCServer()
    b.StartRPCServer()
    
    // This channel will receive one 'true' for each process completion
    complete := make(chan bool)
    
    // Begin processing. If running one coordinator per process, should perhaps
    // run ProcessTurns on the main thread because why the hell not.
    go a.ProcessTurns(complete)
    go b.ProcessTurns(complete)
    
    // Wait for processing to complete
    <- complete
    <- complete
    
    // Yo ho, me hearties, yo ho!
    log.Println("Done")
}
