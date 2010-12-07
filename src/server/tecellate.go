package main

import (
	"fmt"
	"os"
	"json"
	"io/ioutil"
	"log"
	"net"
	"rand"
	"easynet"
	"ttypes"
)

func main() {
	config := loadConfig()
	grid, botConfs := readGridFromFile(os.Args[1])
	connections := connectToCoordinators(config)
	coordConfigs := configureCoordinators(config, botConfs)
	
	for i, conn := range(connections) {
		coordConfigs[i].Terrain = *grid
		data, err := json.Marshal(coordConfigs[i])
		if (err != nil) { log.Exit(err) }
		conn.Write(data)
	}
	
	for i, conn := range(connections) {
		fmt.Printf("Master waiting for first confirmation from %d\n", i)
		fmt.Printf("Master received first confirmation from %d: %s\n", i+1, string(easynet.ReceiveFrom(conn)))
	}
	
	for i, conn := range(connections) {
		fmt.Printf("Master waiting for second confirmation from %d\n", i)
		fmt.Printf("Master received second confirmation from %d: %s\n", i+1, string(easynet.ReceiveFrom(conn)))
	}
	
	fmt.Printf("Starting first turn\n")
	connections[0].Write([]uint8("begin"))
	
	fmt.Println("Master now waiting for results")
	fmt.Printf("Final response from first coordinator: %s\n", 
						 string(easynet.ReceiveFrom(connections[0])))
}

func loadConfig() *ttypes.Config {
	config := new(ttypes.Config)
	
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

func connectToCoordinators(config *ttypes.Config) ([]*net.TCPConn) {
	connections := make([]*net.TCPConn, len(config.Coords))
	for i, _ := range(connections) {
		connections[i] = easynet.Dial(config.Coords[i])
	}
	return connections
}

func configureCoordinators(config *ttypes.Config, botConfs []ttypes.BotConf) ([]ttypes.CoordConfig) {
	coordConfigs := make([]ttypes.CoordConfig, len(config.Coords))
	for i, _ := range(coordConfigs) {
		coordConfigs[i].Identifier = i+1
		coordConfigs[i].NumTurns = config.NumTurns
	}
	for _, conf := range(botConfs) {
		ix := rand.Int() % len(coordConfigs)
		coordConfigs[ix].BotConfs = append(coordConfigs[ix].BotConfs, conf)
	}
	for i, _ := range(coordConfigs) {
		for j, _ := range(coordConfigs) {
			if i != j {
				newAdj := ttypes.AdjacentCoord{coordConfigs[j].Identifier, config.Coords[j]};
				coordConfigs[i].AdjacentCoords = []ttypes.AdjacentCoord{newAdj}
			}
		}
	}
	return coordConfigs
}
