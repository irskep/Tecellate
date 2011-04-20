package util

import (
    "logflow"
    "net"
    "netchan"
    "os"
    "time"
)

func MakeImporterWithRetry(network string, remoteaddr string, n int, log logflow.Logger) *netchan.Importer {
    if log == nil {
        log = logflow.NewSource("util")
    }
    var err os.Error
    for i := 0; i < n; i++ {
        conn, err := net.Dial(network, "", remoteaddr)
        if err == nil {
            return netchan.NewImporter(conn)
        }
        log.Print("Netchan import failed, retrying")
        time.Sleep(1e9/2)
    }
    log.Print("Netchan import failed ", n, " times. Bailing out.")
    log.Fatal(err)
    return nil
}
