package main

import (
	"easynet"
	"fmt"
	"math"
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

func distance(x1 uint, y1 uint, x2 uint, y2 uint) float64 {
	dx := x1-x2
	dy := y1-y2
	return math.Sqrt(float64(dx*dx+dy*dy))
}

func messagesHeardBy(info ttypes.BotInfo) []ttypes.Message {
	messages := make([]ttypes.Message, 0, len(botStates))
	for _, s := range(botStates) {
		lm := s.Info.LastMessage
		d := distance(info.X, info.Y, s.Info.X, s.Info.Y)
		if d <= 2 && len(lm) > 1 {
			messages = append(messages, ttypes.Message{lm, d})
		}
	}
	return messages
}

func moveBots(otherInfos []ttypes.BotInfo) {
	for botNum, s := range(botStates) {
		if s.TurnsToNextMove == 0 {
			req := new(ttypes.BotMoveRequest)
			req.Terrain = config.Terrain
			req.OtherBots = otherInfos
			req.Messages = messagesHeardBy(s.Info)
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
			
			s.Info.LastMessage = rsp.BroadcastMessage
			if len(s.Info.LastMessage) > 1024 {
				s.Info.LastMessage = s.Info.LastMessage[0:1024]
			}
			
			// I could not for the life of me find Go's abs() function.
			s.TurnsToNextMove = oldElevation-newElevation
			if s.TurnsToNextMove < 0 {
				s.TurnsToNextMove = -s.TurnsToNextMove
			}
		} else {
			fmt.Printf("Bot %d hit rocky terrain\n", botNum)
			s.TurnsToNextMove -= 1
			s.Info.LastMessage = ""
		}
	}
}