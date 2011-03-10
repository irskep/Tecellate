package link

import "testing"
import "fmt"

func TestMakeMessage(t *testing.T) {
    fmt.Println(NewMessage(Commands["Look"], "here", "or there", NewMessage(Commands["Listen"])))
}
