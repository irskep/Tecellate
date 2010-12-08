package main

import (
	"ttypes"
)

// All this when a measly filter/map would have done...
func botInfosForNeighbor(neighbor int) []ttypes.BotInfo {
	infos := make([]ttypes.BotInfo, 0, len(botStates))
	
	for _, s := range(botStates) {
		infos = append(infos, s.Info)
	}
	
	return infos
}
