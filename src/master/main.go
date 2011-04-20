/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: master/main.go

MASTER MAIN
*/

package main

import (
    "master"
    "os"
)

func main() {
    m := master.New(os.Args)
    m.ConnectToCoords()
}
