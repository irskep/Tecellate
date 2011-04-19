package testerlib

import "testing"
import "os"
import "fmt"
import "runtime"
import "strings"
import "logflow"

var write_to_sinks logflow.WriteToSinks = logflow.WriteToSinksFunction
func init() {
    runtime.GOMAXPROCS(1)
}

func InitLogs(name string, t *testing.T) (closer func()) {
    fmt.Println("    -", name)

    // Show all output if test fails
    writer := logflow.NewTestWriter(t)
    snk, _ := logflow.NewSink(writer, ".*")

    err := os.MkdirAll("logs/wifi", 0776)
    if err != nil {
        panic("Directory logs/wifi could not be created.")
    }
//     logflow.FileSink("logs/wifi/test/" + name, true, ".*")
//     logflow.StdoutSink("agent/wifi.*")

//     log = func(v ...interface{}) {
//         logflow.Println(fmt.Sprintf("test/%v", name), v...)
//     }

    defer func() {
       logflow.WriteToSinksFunction = func(keypath, s string) {
            if strings.HasPrefix(keypath, "agent/wifi") {
                fmt.Println(keypath, s)
                snk.Write(keypath, s)
            } else if strings.HasPrefix(keypath, fmt.Sprintf("test/%v", name)) {
                snk.Write(keypath, s)
            } else if strings.HasPrefix(keypath, "coord") {
//                 snk.Write(keypath, s)
            }
       }
       writer.Write([]byte(fmt.Sprintf(`

    - Start Testing %v
--------------------------------------------------------------------------------`, name)))
    }()

    closer = func() {
    writer.Write([]byte(fmt.Sprintf(
`--------------------------------------------------------------------------------
    - End Testing %v

`, name)))
    logflow.RemoveAllSinks()
    logflow.WriteToSinksFunction = write_to_sinks
    }
    return closer
}
