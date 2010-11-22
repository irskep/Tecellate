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

type CoordConfig struct {
	Identifier int
	BotConfs []BotConf
	Grid []byte
}

type BotConf struct {
	Path string
}
