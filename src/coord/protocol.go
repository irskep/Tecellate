package main

import (
    "fmt"
    "easynet"
    "json"
    "net"
    "time"
    "ttypes"
)

func listenForMaster(connectionToMaster *net.TCPConn) {
    msg, err := easynet.ReceiveFromWithError(connectionToMaster)
    if err != nil {
        fmt.Printf("%d got an error on the connection to master (because it didn't receive BEGIN, which is fine): %v\n", config.Identifier, err)
    } else {
        if string(msg) == "begin" {
            fmt.Printf("%d is the Chosen One!\n", config.Identifier)
            processing = true
            go processNodes()
        }
    }
}

func listenForPeer() {
    fmt.Printf("%d serving requests\n", config.Identifier)
    completionsRemaining = len(adjsServe)
    neighborsLeftUntilUnlock = len(adjsServe)
    lastInfoListCopied = 0
    for data := range listenServe {
        fmt.Println(string(data))
        //Sometimes requests will be stuck together. Here I am separating them.
        //A crappy and hopefully temporary fix.
        splitPoint := 0
        if data[0] == "{"[0] {
            for i := 1; i < len(data); i++ {
                if data[i-1] == "}"[0] && data[i] == "{"[0] {
                    splitPoint = i
                    break
                }
            }
        } else {
            for i := 1; i < len(data); i++ {
                if data[i] == "{"[0] {
                    splitPoint = i
                    break
                }
            }
        }
        if splitPoint == 0 {
            handleRequest(data)
        } else {
            fmt.Println("Split occurred, wish I knew how to flush buffers...")
            handleRequest(data[0:splitPoint])
            handleRequest(data[splitPoint:len(data)])
        }
    }
}

func handleRequest(data []uint8) {
    if processing == false {
        // The game's afoot!
        processing = true
        go processNodes()
    }
    r := new(Request)
    err := json.Unmarshal(data, r)
    easynet.DieIfError(err, "JSON error")
    switch {
    case r.Command == "GetNodes":
        fmt.Printf("%d handle GetNodes from %d\n", config.Identifier, r.Identifier)
        for respondingToRequestsFor < r.Turn {
            fmt.Printf("%d not ready for GetNodes\n", config.Identifier)
            time.Sleep(1000000)
        }

        fmt.Printf("%d ready for GetNodes\n", config.Identifier)
        info := new(RespondNodeInfo)
        info.Identifier = config.Identifier
        info.Turn = respondingToRequestsFor
        if lastInfoListCopied == 0 {
            info.BotData = botInfosForNeighbor(r.Identifier)
        } else {
            info.BotData = infoQueue[lastInfoListCopied][0:len(botStates)]
        }
        easynet.SendJson(adjsServe[r.Identifier], info)
        fmt.Printf("%d sent GetNodes response to %d\n", config.Identifier, r.Identifier)
        fmt.Printf("    and it was %v\n", info)

        neighborsLeftUntilUnlock -= 1
        if neighborsLeftUntilUnlock <= 0 {
            fmt.Printf("served all requests for %d\n", lastInfoListCopied)
            lastInfoListCopied += 1
            neighborsLeftUntilUnlock = len(adjsServe)
        }
    case r.Command == "Complete":
        completionsRemaining -= 1
        if completionsRemaining == 0 {
            fmt.Println("All neighbors complete, signaling TCOMPLETE2")
            complete <- true
        }
    }
}

func processNodes() {
    fmt.Printf("%d processing nodes\n", config.Identifier)

    for respondingToRequestsFor < config.NumTurns {
        for lastInfoListCopied < respondingToRequestsFor {
            fmt.Printf("%d wants to process %d but waits for %d\n", config.Identifier, respondingToRequestsFor, lastInfoListCopied)
            time.Sleep(1000000)
        }
        fmt.Printf("%d starting turn %d\n", config.Identifier, respondingToRequestsFor)

        if respondingToRequestsFor > 0 {
            fmt.Println(infoQueue)
            for k, s := range botStates {
                s.Info = infoQueue[respondingToRequestsFor][k]
            }
        }

        otherInfos := getAgentInfoFromNeighbors()

        fmt.Printf("EVERYONE SHOULD SEE %v\n", otherInfos)

        declareDeaths(otherInfos)

        moveBots(otherInfos)

        respondingToRequestsFor += 1
        infoQueue = append(infoQueue, otherInfos)
    }
    broadcastComplete()
    fmt.Println("Turns complete, signaling TCOMPLETE1")
    complete <- true
}

func getAgentInfoFromNeighbors() []ttypes.BotInfo {
    otherInfos := make([]ttypes.BotInfo, len(botStates), len(botStates)*len(adjsServe))

    //Copy all infos from botStates into otherInfos
    for k, s := range botStates {
        otherInfos[k] = s.Info
    }

    //Get updates from neighbors
    for j, conn := range adjsRequest {
        fmt.Printf("%d turn %d, request neighbor %d\n", config.Identifier, respondingToRequestsFor, j)
        r := new(Request)
        r.Identifier = config.Identifier
        r.Turn = respondingToRequestsFor
        r.Command = "GetNodes"

        easynet.SendJson(conn, r)

        info := new(RespondNodeInfo)
        easynet.ReceiveJson(conn, info)

        otherInfos = append(otherInfos, info.BotData...)
    }
    return otherInfos
}

func broadcastComplete() {
    note := new(Request)
    note.Identifier = config.Identifier
    note.Turn = respondingToRequestsFor
    note.Command = "Complete"

    for i, conn := range adjsRequest {
        fmt.Printf("%d broadcasting complete to %d\n", config.Identifier, i)
        easynet.SendJson(conn, note)
    }
}
