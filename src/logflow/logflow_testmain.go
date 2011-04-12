package main

import "logflow"

import (
    "fmt"
    "testing"
    __regexp__ "regexp"
)

var test_logflow = []testing.InternalTest {
	testing.InternalTest{ "logflow.TestInstantiate", logflow.TestInstantiate },
	testing.InternalTest{ "logflow.TestNothing", logflow.TestNothing },
}


func main() {
	fmt.Println("Testing logflow:");
	testing.Main(__regexp__.MatchString, test_logflow);
}
