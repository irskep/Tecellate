package link

import "testing"
import "fmt"
import . "byteslice"

func TestMakeMessage(t *testing.T) {
    fmt.Println(NewMessage(Commands["Look"], ByteSlice([]byte("here")), ByteSlice([]byte("or there"))))
}
