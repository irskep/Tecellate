package main

import (
	"fmt"
	"os"
	"json"
	"io/ioutil"
	"log"
	"net"
	"rand"
)

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
	identifier int
	BotConfs []BotConf
}

type BotConf struct {
	Path string
}

func main() {
	config := loadConfig()
	connections := connectToCoordinators(config)
	coordConfigs := configureCoordinators(config)
	for i, conn := range(connections) {
		data, err := json.Marshal(coordConfigs[i])
		if (err != nil) { log.Exit(err) }
		conn.Write(data)
	}
}

func loadConfig() *Config {
	config := new(Config)
	
	configFile, err := os.Open("config.json", os.O_RDONLY, 0)
	if err != nil { log.Exit(err) }
	defer configFile.Close()
	
	configBytes, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Exit(err)
	} else {
		json.Unmarshal(configBytes, config)
	}
	return config
}

func connectToCoordinators(config *Config) ([]*net.TCPConn) {
	connections := make([]*net.TCPConn, len(config.Coords))
	for i, _ := range(connections) {
		fmt.Printf("%d", i)
		addr, err := net.ResolveTCPAddr(config.Coords[i]);
		if err != nil { log.Exit(err) }
		conn, err := net.DialTCP("tcp", nil, addr)
		if err != nil { log.Exit(err) }
		connections[i] = conn
	}
	return connections
}

func configureCoordinators(config *Config) ([]CoordConfig) {
	coordConfigs := make([]CoordConfig, len(config.Coords))
	for _, bot := range(config.BotDefs) {
		for i := 0; i < bot.Count; i++ {
			ix := rand.Int() % len(coordConfigs)
			newConf := BotConf{bot.Path}
			coordConfigs[ix].BotConfs = append(coordConfigs[ix].BotConfs, newConf)
		}
	}
	return coordConfigs
}
