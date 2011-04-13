package main

import "logflow"

import (
    "fmt"
    "testing"
    __regexp__ "regexp"
)

var test_logflow = []testing.InternalTest {
    testing.InternalTest{ "logflow.TestHookup", logflow.TestHookup },
}


func main() {
	fmt.Println("Testing logflow:");
	testing.Main(__regexp__.MatchString, test_logflow);
}
