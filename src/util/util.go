package util

import (
    "logflow"
    "net"
    "netchan"
    "time"
)

func MakeImporterWithRetry(network string, remoteaddr string, n int) *netchan.Importer {
    var err os.Error
    for i := 0; i < n; i++ {
        conn, err := net.Dial(network, "", remoteaddr)
        if err == nil {
            return netchan.NewImporter(conn)
        }
        logflow.Print("agent", "Netchan import failed, retrying")
        time.Sleep(1e9/2)
    }
    logflow.Print("agent", "Netchan import failed ", n, " times. Bailing out.")
    logflow.Fatal("agent", err)
    return nil
}
