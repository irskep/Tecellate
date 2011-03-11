package agent

type Agent interface {
    Turn() bool
    State() *AgentState
    ApplyTransform(Transform)
}
