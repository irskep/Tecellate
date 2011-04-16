package wifi

import "testing"
import "os"
import "fmt"
import "runtime"
import "strings"
import "logflow"

var write_to_sinks logflow.WriteToSinks = logflow.WriteToSinksFunction
func init() {
    runtime.GOMAXPROCS(4)
}

func initLogs(name string, t *testing.T) func() {
    // Show all output if test fails
    logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs/wifi", 0776)
    if err != nil {
        panic("Directory logs/wifi could not be created.")
    }
//     logflow.FileSink("logs/wifi/test/" + name, true, ".*")
//     logflow.StdoutSink("agent/wifi.*")

    defer func() {
       logflow.WriteToSinksFunction = func(keypath, s string) {
           if strings.HasPrefix(keypath, "agent/wifi") {
               fmt.Print(keypath, ": ", s)
           }
       }
    }()

    defer logflow.Println("test", fmt.Sprintf(`
--------------------------------------------------------------------------------
    Start Testing %v
`, name))
    return func() {
    logflow.Println("test", fmt.Sprintf(`
--------------------------------------------------------------------------------
    End Testing %v
`, name))
    logflow.RemoveAllSinks()
    logflow.WriteToSinksFunction = write_to_sinks
    }
}
