package logflow

import (
    "bytes"
    "fmt"
    "os"
)

/*
    Types and stuff
*/

type LogLevel string
const (
    INFO LogLevel = "info"
    DEBUG = "debug"
    WARN = "warn"
    ERROR = "error"
    FATAL = "fatal"
)

type Source interface {
    Log(LogLevel, ...interface{})
    Logf(LogLevel, string, ...interface{})
    Logln(LogLevel, ...interface{})
    LOutput(LogLevel, int, string) os.Error
}

type Logger interface {
    Fatal(...interface{})
    Fatalf(string, ...interface{})
    Fatalln(...interface{})
    Output(int, string) os.Error
    Panic(...interface{})
    Panicf(string, ...interface{})
    Panicln(...interface{})
    Print(...interface{})
    Printf(string, ...interface{})
    Println(...interface{})
    Log(LogLevel, ...interface{})
    Logf(LogLevel, string, ...interface{})
    Logln(LogLevel, ...interface{})
}

type source struct {
    keypath string
    sinks []Sink
}

/*
    Things that make me a special snowflake
*/

func NewSource(keypath string) *source {
    theSource := &source{keypath: keypath, sinks: SinksMatchingKeypath(keypath)}
    return theSource
}

func (self *source) LOutput(level LogLevel, s string) os.Error {
    buf := new(bytes.Buffer)
    buf.WriteString(self.keypath)
    buf.WriteByte('/')
    buf.WriteString(string(level))
    buf.WriteString(": ")
    buf.WriteString(s)
    if len(s) > 0 && s[len(s)-1] != '\n' {
        buf.WriteByte('\n')
    }
    WriteToSinksMatchingKeypath(fmt.Sprintf("%v/%v", self.keypath, level), buf.Bytes())
    return nil
}

func (self *source) String() string {
    return fmt.Sprintf("<Source: %v>", self.keypath)
}

/*
    Implement Logging interface
*/

func (self *source) Fatal(v ...interface{}) {
    self.LOutput(FATAL, fmt.Sprint(v...))
    os.Exit(1)
}

func (self *source) Fatalf(format string, v ...interface{}) {
    self.LOutput(FATAL, fmt.Sprintf(format, v...))
    os.Exit(1)
}

func (self *source) Fatalln(v ...interface{}) {
    self.LOutput(FATAL, fmt.Sprintln(v...))
    os.Exit(1)
}

func (self *source) Output(calldepth int, s string) os.Error {
    // Discard calldepth
    return self.LOutput(INFO, s)
}

func (self *source) Panic(v ...interface{}) {
    s := fmt.Sprint(v...)
    self.LOutput(ERROR, s)
    panic(s)
}

func (self *source) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    self.LOutput(ERROR, s)
    panic(s)
}

func (self *source) Panicln(v ...interface{}) {
    s := fmt.Sprintln(v...)
    self.LOutput(ERROR, s)
    panic(s)
}

func (self *source) Print(v ...interface{}) {
    self.LOutput(INFO, fmt.Sprint(v...))
}

func (self *source) Printf(format string, v ...interface{}) {
    self.LOutput(INFO, fmt.Sprintf(format, v...))
}

func (self *source) Println(v ...interface{}) {
    self.LOutput(INFO, fmt.Sprintln(v...))
}

func (self *source) Log(level LogLevel, v ...interface{}) {
    self.LOutput(level, fmt.Sprint(v...))
}

func (self *source) Logf(level LogLevel, format string, v ...interface{}) {
    self.LOutput(level, fmt.Sprintf(format, v...))
}

func (self *source) Logln(level LogLevel, v ...interface{}) {
    self.LOutput(level, fmt.Sprintln(v...))
}
