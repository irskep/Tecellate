package main

import (
    "net"
    "ttypes"
)

type Request struct {
    Identifier int
    Turn       int
    Command    string
}

type RespondNodeInfo struct {
    Identifier int
    Turn       int
    BotData    []ttypes.BotInfo
}

type BotState struct {
    Conn *net.TCPConn
    Info ttypes.BotInfo
}

func (s BotState) Dead() bool {
    return s.Info.Killed == true
}

func (s BotState) CollidesWith(i ttypes.BotInfo) bool {
    return s.Info.X == i.X && s.Info.Y == i.Y
}
