package logflow

import (
    "io"
    "os"
    "regexp"
    "strings"
)

type ClosingWriter interface {
    io.Writer
    io.Closer
}

type Sink interface {
    MatchesKeypath(string) bool
    Write([]byte)
}

var sinks []*sink = make([]*sink, 0)

type sink struct {
    keypathRegexp *regexp.Regexp
    writer ClosingWriter
}

func NewSink(w ClosingWriter, matches ...string) (*sink, os.Error) {
    re, err := regexp.Compile(strings.Join([]string{"^", strings.Join(matches, "|"), "$"}, ""))
    var theSink *sink = nil
    if err == nil {
        theSink = &sink{keypathRegexp: re, writer: w}
        sinks = append(sinks, theSink)
    }
    return theSink, err
}

func StdoutSink(matches ...string) (*sink, os.Error) {
    return NewSink(os.Stdout, matches...)
}

func StderrSink(matches ...string) (*sink, os.Error) {
    return NewSink(os.Stderr, matches...)
}

func FileSink(path string, matches ...string) (*sink, os.Error) {
    f, err := os.Open("logs/agents", os.O_RDWR|os.O_CREAT, 0664)
    if err != nil {
        return nil, err
    }
    return NewSink(f, "agent/.*")
}

func SinksMatchingKeypath(keypath string) []Sink {
    matches := make([]Sink, 0)
    for _, snk := range sinks {
        if snk.MatchesKeypath(keypath) {
            matches = append(matches, snk)
        }
    }
    return matches
}

func WriteToSinksMatchingKeypath(keypath string, s []byte) {
    for _, snk := range sinks {
        if snk.MatchesKeypath(keypath) {
            snk.Write(s)
        }
    }
}

func RemoveAllSinks() {
    for _, snk := range sinks {
        snk.writer.Close()
    }
    sinks = make([]*sink, 0)
}

func (self *sink) String() string {
    return self.keypathRegexp.String()
}

func (self *sink) MatchesKeypath(keypath string) bool {
    return self.keypathRegexp.MatchString(keypath)
}

func (self *sink) Write(s []byte) {
    self.writer.Write(s)
}
