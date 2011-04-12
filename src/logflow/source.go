package logflow

import (
    "fmt"
    // "log"
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

type source struct {
    keypath string
    sinks []Sink
}

/*
    Things that make me a special snowflake
*/

func NewSource(keypath string) *source {
    return &source{
        keypath: keypath,
        sinks: nil,
    }
}

func (self *source) LOutput(level LogLevel, calldepth int, s string) os.Error {
    return nil
}

func (self *source) String() string {
    return fmt.Sprintf("<Source: %v>", self.keypath)
}

/*
    Implement Logging interface
*/

func (self *source) Fatal(v ...interface{}) {
    self.LOutput(FATAL, 2, fmt.Sprint(v...))
    os.Exit(1)
}

func (self *source) Fatalf(format string, v ...interface{}) {
    self.LOutput(FATAL, 2, fmt.Sprintf(format, v...))
    os.Exit(1)
}

func (self *source) Fatalln(v ...interface{}) {
    self.LOutput(FATAL, 2, fmt.Sprintln(v...))
    os.Exit(1)
}

func (self *source) Output(calldepth int, s string) os.Error {
    return self.LOutput(INFO, calldepth, s)
}

func (self *source) Panic(v ...interface{}) {
    s := fmt.Sprint(v...)
    self.LOutput(ERROR, 2, s)
    panic(s)
}

func (self *source) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    self.LOutput(ERROR, 2, s)
    panic(s)
}

func (self *source) Panicln(v ...interface{}) {
    s := fmt.Sprintln(v...)
    self.LOutput(ERROR, 2, s)
    panic(s)
}

func (self *source) Print(v ...interface{}) {
    self.LOutput(INFO, 2, fmt.Sprint(v...))
}

func (self *source) Printf(format string, v ...interface{}) {
    self.LOutput(INFO, 2, fmt.Sprintf(format, v...))
}

func (self *source) Println(v ...interface{}) {
    self.LOutput(INFO, 2, fmt.Sprintln(v...))
}

func (self *source) Log(level LogLevel, v ...interface{}) {
    self.LOutput(level, 2, fmt.Sprint(v...))
}

func (self *source) Logf(level LogLevel, format string, v ...interface{}) {
    self.LOutput(level, 2, fmt.Sprintf(format, v...))
}

func (self *source) Logln(level LogLevel, v ...interface{}) {
    self.LOutput(level, 2, fmt.Sprintln(v...))
}
