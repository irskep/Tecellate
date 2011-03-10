package link

import "fmt"
import "strings"

const Timeout = 1e9

type Command uint8
var Commands map[string]Command
var cmdsr []string

type Argument interface{}
type Arguments []Argument

func init() {
    Commands = make(map[string]Command)
    cmdsr = []string{
        "Ack", "Nak", "Move", "Look", "Collect", "Listen", "Broadcast",
        "Complete", "Start", "Exit",
    }
    for i, cmd := range cmdsr {
        Commands[cmd] = Command(i)
    }
}

type Message struct {
    Cmd Command
    Args Arguments
}

func NewMessage(cmd Command, args ... Argument) *Message {
    return &Message{Cmd: cmd, Args: args}
}

func (self Message) String() string {
    return fmt.Sprintf("<Message cmd:%s args:%s>", cmdsr[self.Cmd], self.Args)
}

func (self Arguments) String() string {
    s := make([]string, 0, 3)
    args := make([]string, 0, len(self))
    s = append(s, "[")
    for _, arg := range self {
        args = append(args, fmt.Sprintf("{%s}", arg))
    }
    s = append(s, strings.Join(args, ", "))
    s = append(s, "]")
    return strings.Join(s, "")
}

type Link chan Message
