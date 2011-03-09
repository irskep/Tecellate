package game

type Config struct {
    AgentStarts []AgentStart
}

type AgentStart struct {
    Position Point
    Kind string
}
