package agent

type Agent interface {
    Turn() bool
    State() *AgentState
    Apply(Transform)
    MigrateTo(string)
}
