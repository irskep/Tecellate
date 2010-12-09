package main

import (
	"net"
	"ttypes"
)

type BotState struct {
	Conn *net.TCPConn
	Info ttypes.BotInfo
	TurnsToNextMove int
	Killed bool
}

type Request struct {
	Identifier int
	Turn int
	Command string
}

type RespondNodeInfo struct {
	Identifier int
	Turn int
	BotData []ttypes.BotInfo
}