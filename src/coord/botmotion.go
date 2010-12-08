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
		if s.Killed == false {
			infos = append(infos, s.Info)
		}
	}
	
	return infos
}

func distance(x1 uint, y1 uint, x2 uint, y2 uint) float64 {
	dx := x1-x2
	dy := y1-y2
	return math.Sqrt(float64(dx*dx+dy*dy))
}

func perceptionOf(info ttypes.BotInfo) ([]ttypes.Message, []ttypes.BotInfo) {
	messages := make([]ttypes.Message, 0, len(botStates))
	otherBots := make([]ttypes.BotInfo, 0, 30)
	for _, s := range(botStates) {
		if s.Killed == false {
			lm := s.Info.LastMessage
			d := distance(info.X, info.Y, s.Info.X, s.Info.Y)
			if d <= 5 && len(lm) > 1 {
				messages = append(messages, ttypes.Message{lm, d})
			}
			if d <= 3 {
				otherBots = append(otherBots, s.Info)
			}
		}
	}
	return messages, otherBots
}

// Stupid n^2 algorithm to see if any 2 bots overlap and mark them killed if they do
func declareDeaths(otherInfos []ttypes.BotInfo) {
	for ix, s := range(botStates) {
		if s.Killed == false {
			for jx, oi := range(otherInfos) {
				if ix != jx && s.Info.X == oi.X && s.Info.Y == oi.Y {
					fmt.Printf("Killing bot %v\n", s)
					s.Killed = true
				}
			}
		}
	}
	fmt.Println(botStates[0].Killed)
}

func moveBots(otherInfos []ttypes.BotInfo) {
	// fmt.Printf("All infos: %v\n", otherInfos)
	for botNum, s := range(botStates) {
		if s.Killed == false {
			if s.TurnsToNextMove == 0 {
				msges, botsSeen := perceptionOf(s.Info)
				req := new(ttypes.BotMoveRequest)
				req.Terrain = config.Terrain
				req.OtherBots = botsSeen
				req.Messages = msges
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
			
				otherInfos[botNum].LastMessage = rsp.BroadcastMessage
				if len(otherInfos[botNum].LastMessage) > 1024 {
					otherInfos[botNum].LastMessage = otherInfos[botNum].LastMessage[0:1024]
				}
			
				// I could not for the life of me find Go's abs() function.
				s.TurnsToNextMove = oldElevation-newElevation
				if s.TurnsToNextMove < 0 {
					s.TurnsToNextMove = -s.TurnsToNextMove
				}
			} else {
				fmt.Printf("Bot %d hit rocky terrain\n", botNum)
				s.TurnsToNextMove -= 1
				otherInfos[botNum].LastMessage = ""
			}
		}
	}
}