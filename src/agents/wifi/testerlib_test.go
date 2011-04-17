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

func initLogs(name string, t *testing.T) (log func(...interface{}), closer func()) {
    fmt.Println("    -", name)

    // Show all output if test fails
    snk, _ := logflow.NewSink(logflow.NewTestWriter(t), ".*")

    err := os.MkdirAll("logs/wifi", 0776)
    if err != nil {
        panic("Directory logs/wifi could not be created.")
    }
//     logflow.FileSink("logs/wifi/test/" + name, true, ".*")
//     logflow.StdoutSink("agent/wifi.*")

    log = func(v ...interface{}) {
        logflow.Println(fmt.Sprintf("test/%v", name), v...)
    }

    defer func() {
       logflow.WriteToSinksFunction = func(keypath, s string) {
           if strings.HasPrefix(keypath, "agent/wifi") {
               snk.Write(keypath, s)
           } else if strings.HasPrefix(keypath, fmt.Sprintf("test/%v", name)) {
               snk.Write(keypath, s)
           }
       }
       log(fmt.Sprintf(`
--------------------------------------------------------------------------------
    Start Testing %v
`, name))
    }()

    closer = func() {
    log(fmt.Sprintf(`
--------------------------------------------------------------------------------
    End Testing %v
`, name))
    logflow.RemoveAllSinks()
    logflow.WriteToSinksFunction = write_to_sinks
    }
    return log, closer
}
