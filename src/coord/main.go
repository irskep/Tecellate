/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/main.go

COORD MAIN
*/

package main

import (
    "coord/runner"
    "logflow"
    "os"
)

func main() {
    logflow.StdoutSink(".*")
    runner.RunAtAddress(os.Args[1])
}
