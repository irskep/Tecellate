package link

import "fmt"
import "strings"
import . "byteslice"

const Timeout = 10e9

type SendLink chan<- Message
type RecvLink <-chan Message

type Command uint8
var Commands map[string]Command
var cmdsr []string

type Arguments []ByteSlice

type Argument interface {
    Bytes() ByteSlice
}

func init() {
    Commands = make(map[string]Command)
    cmdsr = []string{
        "Ack", "Nak", "Move", "Look", "Collect", "Listen", "Broadcast",
        "Complete", "Start", "Exit", "PrevResult", "Id", "Energy", "Migrate", 
    }
    for i, cmd := range cmdsr {
        Commands[cmd] = Command(i)
    }
}

type Message struct {
    Cmd Command
    Args Arguments
}

func (self Command) Bytes() ByteSlice {
    return ByteSlice8(uint8(self))
}

func MakeCommand(bytes ByteSlice) Command {
    return Command(bytes.Int8())
}

func NewMessage(cmd Command, args ... Argument) *Message {
    arguments := make(Arguments, 0, len(args))
    for _, arg := range args {
        arguments = append(arguments, arg.Bytes())
    }
    return &Message{Cmd: cmd, Args: arguments}
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
