package logflow

import (
    "fmt"
    // "log"
    "os"
)

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
}

type source struct {
    keypath string
    sinks []Sink
}

func NewSource(keypath string) *source {
    return &source{
        keypath: keypath,
        sinks: nil,
    }
}

func (self *source) String() string {
    return fmt.Sprintf("%v", self.keypath)
}

func (self *source) Fatal(v ...interface{}) {
    self.Output(2, fmt.Sprint(v...))
    os.Exit(1)
}

func (self *source) Fatalf(format string, v ...interface{}) {
    self.Output(2, fmt.Sprintf(format, v...))
    os.Exit(1)
}

func (self *source) Fatalln(v ...interface{}) {
    self.Output(2, fmt.Sprintln(v...))
    os.Exit(1)
}

func (self *source) Output(calldepth int, s string) os.Error {
    return nil
}

func (self *source) Panic(v ...interface{}) {
    s := fmt.Sprint(v...)
    self.Output(2, s)
    panic(s)
}

func (self *source) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    self.Output(2, s)
    panic(s)
}

func (self *source) Panicln(v ...interface{}) {
    s := fmt.Sprintln(v...)
    self.Output(2, s)
    panic(s)
}

func (self *source) Print(v ...interface{}) {
    self.Output(2, fmt.Sprint(v...))
}

func (self *source) Printf(format string, v ...interface{}) {
    self.Output(2, fmt.Sprintf(format, v...))
}

func (self *source) Println(v ...interface{}) {
    self.Output(2, fmt.Sprintln(v...))
}

func (self *source) Log(level LogLevel, v ...interface{}) {
    
}

func (self *source) Logf(level LogLevel, format string, v ...interface{}) {
    
}

func (self *source) Logln(level LogLevel, v ...interface{}) {
    
}
