package link

import "fmt"
import "strings"

const Timeout = 10e9

type SendLink chan<- Message
type RecvLink <-chan Message

type Command uint8
var Commands map[string]uint8
var cmdsr []string

type Argument interface{}
type Arguments []Argument

func init() {
    Commands = make(map[string]uint8)
    cmdsr = []string{
        "Ack", "Nak", "Move", "Look", "Collect", "Listen", "Broadcast",
        "Complete", "Start", "Exit", "PrevResult", "Id", "Energy",
    }
    for i, cmd := range cmdsr {
        Commands[cmd] = uint8(i)
    }
}

type Message struct {
    Cmd uint8
    Args Arguments
}

func NewMessage(cmd uint8, args ... Argument) *Message {
    return &Message{Cmd: cmd, Args: args}
}

func (self Message) String() string {
    return fmt.Sprintf("<Message cmd:%s args:%s>", self.Cmd, self.Args)
}

func (self Arguments) String() string {
    s := make([]string, 0, 3)
    args := make([]string, 0, len(self))
    s = append(s, "[")
    for _, arg := range self {
        args = append(args, fmt.Sprintf("{%v}", arg))
    }
    s = append(s, strings.Join(args, ", "))
    s = append(s, "]")
    return strings.Join(s, "")
}

func (self Command) String() string {
    return cmdsr[self]
}
