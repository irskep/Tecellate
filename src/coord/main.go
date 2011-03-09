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
)

func main() {
    fmt.Println("coord/main.go")
    
    a := coord.NewCoordinator()
    b := coord.NewCoordinator()
    
    a.ConnectToLocal(b)
    b.ConnectToLocal(a)
}
