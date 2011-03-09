package coord

import "coord.game"

type CoordinatorProxy struct {
    conn chan []byte
}

func (self *CoordinatorProxy) RequestStatesInBox(bottomLeft Point, topRight Point, turn int) []AgentState {
    return nil;
}

func (self *CoordinatorProxy) SendComplete() {
}
