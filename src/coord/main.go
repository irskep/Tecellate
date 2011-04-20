/*
Tecellate
Authors: Tim Henderson      &    Stephen Johnson
Contact: tim.tadh@gmail.com &    steve@steveasleep.com
File: coord/main.go

COORD MAIN
*/

package main

import (
    "coord"
    "coord/runner"
    "os"
)

func main() {
    c := coord.NewCoordinator()
    r := runner.New(c, os.Args[1])
    r.ExportNetchans()
    r.RunExporter()
    r.ReadConfig()
}
