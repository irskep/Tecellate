package logflow

import "testing"
import "fmt"

func TestSourceInstantiate(t *testing.T) {
    src := NewSource("test.info")
    fmt.Println(src)
}

func TestSinkInstantiate(t *testing.T) {
    sink, err := NewSink("test/info", "test/debug")
    fmt.Println(sink, err)
}