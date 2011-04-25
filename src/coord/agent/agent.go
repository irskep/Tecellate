package agent

type Agent interface {
    Stop()
    Turn() bool
    State() *AgentState
    Apply(Transform)
    MigrateTo(string)
}
