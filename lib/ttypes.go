package ttypes

import (
	"strconv"
	"strings"
)

type Config struct {
	Coords []string
	NumTurns int
	SplitStrategy string
}

type AdjacentCoord struct {
	Identifier int
	Address string
}

type CoordConfig struct {
	Identifier int
	NumTurns int
	BotConfs []BotConf
	Terrain Grid
	AdjacentCoords []AdjacentCoord
}

type BotConf struct {
	Path string
	X uint
	Y uint
}

type BotInfo struct {
	X uint
	Y uint
	LastMessage string
}

type Message struct {
	Body string
	Distance float64
}

type BotMoveRequest struct {
	Terrain Grid
	OtherBots []BotInfo
	Messages []Message
	YourX uint
	YourY uint
	Kill bool
}

type BotMoveResponse struct {
	MoveDirection string
	BroadcastMessage string
}

type Grid struct {
	Items []int
	Width uint
	Height uint
}

func (g Grid) Get(x uint, y uint) int {
	return g.Items[x*g.Width+y]
}

func (g Grid) Set(x uint, y uint, val int) {
	g.Items[x*g.Width+y] = val
}

func (g Grid) String() string {
	rowStrings := make([]string, g.Height+1)
	rowStrings[0] = strings.Join([]string{strconv.Uitoa(g.Width), strconv.Uitoa(g.Height)}, " ")
	for y := uint(0); y < g.Height; y++ {
		row := make([]string, g.Width)
		for x := uint(0); x < g.Width; x++ {
			row[x] = strconv.Uitoa(uint(g.Get(x, y)))
		}
		rowStrings[y+1] = strings.Join(row, " ")
	}
	return strings.Join(rowStrings, "\n")
}
