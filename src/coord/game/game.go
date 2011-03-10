package game

import "coord/agent"
import "coord/config"

type GameState struct {
    Turn int
    Agents []*agent.Agent
}

func NewGameState() *GameState {
    return &GameState{0, make([]*agent.Agent, 0)}
}

func (self *GameState) Configure(conf config.Config) {

}

func (self *GameState) ApplyMoves(moves []*agent.Move, agentStates []*agent.AgentState) {

}
