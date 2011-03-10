package agent

type AgentProxy struct {
    State AgentState
    conn chan []byte
}

func (self *AgentProxy) Turn() {

}
