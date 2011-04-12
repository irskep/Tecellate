package logflow

import "testing"
import "fmt"

func TestInstantiate(t *testing.T) {
    src := NewSource("test.basic")
    fmt.Println(src)
}

func TestNothing(t *testing.T) {
    
}