package main

import "logflow"

import (
    "fmt"
    "testing"
    __regexp__ "regexp"
)

var test_logflow = []testing.InternalTest {
	testing.InternalTest{ "logflow.TestSourceInstantiate", logflow.TestSourceInstantiate },
	testing.InternalTest{ "logflow.TestSinkInstantiate", logflow.TestSinkInstantiate },
}


func main() {
	fmt.Println("Testing logflow:");
	testing.Main(__regexp__.MatchString, test_logflow);
}
