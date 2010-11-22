package ttypes

type Config struct {
	Coords []string
	NumCoords int
	ListenAddr string
	BotDefs []BotDef
}

type BotDef struct {
	Path string
	Count int
}

type Grid struct {
	Items []byte
	Width uint
	Height uint
}

type AdjacentCoord struct {
	Identifier int
	Address string
}

type CoordConfig struct {
	Identifier int
	BotConfs []BotConf
	Terrain Grid
	AdjacentCoords []AdjacentCoord
}

type BotConf struct {
	Path string
}
