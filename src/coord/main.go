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
    coords := coord.CoordinatorList(3, &config.Config{0, nil, "none", false, true})
    
    coord.ConnectInChain(coords)
    
    // This channel will receive one 'true' for each process completion
    complete := make(chan bool)
    
    for _, c := range(coords) {
        c.StartRPCServer()
    }
    
    for _, c := range(coords) {
        go c.ProcessTurns(complete)
    }
    
    // Wait for processing to complete
    for _, _ = range(coords) {
        <- complete
    }
    
    // Yo ho, me hearties, yo ho!
    log.Println("Done")
}
