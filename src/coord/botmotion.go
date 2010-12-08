package main

import (
	"easynet"
	"fmt"
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

func moveBots(otherInfos []ttypes.BotInfo) {
	for botNum, s := range(botStates) {
		if s.TurnsToNextMove == 0 {
			req := new(ttypes.BotMoveRequest)
			req.Terrain = config.Terrain
			req.OtherBots = otherInfos
			req.Messages = nil
			req.YourX = s.Info.X
			req.YourY = s.Info.Y
			req.Kill = false
			easynet.SendJson(s.Conn, req)

			rsp := new(ttypes.BotMoveResponse)
			easynet.ReceiveJson(s.Conn, rsp)
			
			oldElevation := config.Terrain.Get(s.Info.X, s.Info.Y)
			
			switch {
			case rsp.MoveDirection == "left" && otherInfos[botNum].X > 0:
				otherInfos[botNum].X -= 1
			case rsp.MoveDirection == "right" && otherInfos[botNum].X < config.Terrain.Width:
				otherInfos[botNum].X += 1
			case rsp.MoveDirection == "down" && otherInfos[botNum].Y > 0:
				otherInfos[botNum].Y -= 1
			case rsp.MoveDirection == "up" && otherInfos[botNum].Y < config.Terrain.Height:
				otherInfos[botNum].Y += 1
			}
			newElevation := config.Terrain.Get(otherInfos[botNum].X, otherInfos[botNum].Y)
			
			// I could not for the life of me find Go's abs() function.
			botStates[botNum].TurnsToNextMove = oldElevation-newElevation
			if botStates[botNum].TurnsToNextMove < 0 {
				botStates[botNum].TurnsToNextMove = -botStates[botNum].TurnsToNextMove
			}
		} else {
			fmt.Printf("Bot %d hit rocky terrain\n", botNum)
			botStates[botNum].TurnsToNextMove -= 1
		}
	}
}