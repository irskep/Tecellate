package ttypes

type Config struct {
	Coords []string
	NumTurns int
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
}

type Message struct {
	Body string
	SourceX uint
	SourceY uint
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
	Items []byte
	Width uint
	Height uint
}

func (g Grid) Get(x uint, y uint) byte {
	return g.Items[x*g.Width+y]
}

func (g Grid) Set(x uint, y uint, val byte) {
	g.Items[x*g.Width+y] = val
}