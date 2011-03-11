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
    coords := coord.ChainedLocalCoordinators(3, config.BasicTestConfig())
    coords.Run()
    
    // Yo ho, me hearties, yo ho!
    log.Println("Done")
}
