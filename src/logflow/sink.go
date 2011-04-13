package logflow

import (
    "bytes"
    "fmt"
    "io"
    "os"
    "regexp"
    "strings"
    "sync"
)

type Sink interface {
    MatchesKeypath(string) bool
    Write(string, string)
    SetWritesPrefix(bool)
}

var sinks []*sink = make([]*sink, 0)

type sink struct {
    keypathRegexp *regexp.Regexp
    writer io.WriteCloser
    writesPrefix bool
    needsClose bool
    mu sync.Mutex
}

func NewSink(w io.WriteCloser, matches ...string) (*sink, os.Error) {
    re, err := regexp.Compile(strings.Join([]string{"^", strings.Join(matches, "|"), "$"}, ""))
    var theSink *sink = nil
    if err == nil {
        theSink = &sink{keypathRegexp: re, writer: w, writesPrefix: true, needsClose: true}
        sinks = append(sinks, theSink)
    }
    return theSink, err
}

func StdoutSink(matches ...string) (*sink, os.Error) {
    newSink, err := NewSink(os.Stdout, matches...)
    if newSink == nil {
        return nil, err
    }
    newSink.needsClose = false
    return newSink, err
}

func StderrSink(matches ...string) (*sink, os.Error) {
    newSink, err := NewSink(os.Stderr, matches...)
    if newSink == nil {
        return nil, err
    }
    newSink.needsClose = false
    return newSink, err
}

func FileSink(path string, matches ...string) (*sink, os.Error) {
    f, err := os.Open(path, os.O_RDWR|os.O_CREAT, 0664)
    if err != nil {
        return nil, err
    }
    newSink, err := NewSink(f, matches...)
    if newSink == nil {
        return nil, err
    }
    newSink.needsClose = true
    return newSink, err
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

func WriteToSinksMatchingKeypath(keypath string, s string) {
    for _, snk := range sinks {
        if snk.MatchesKeypath(keypath) {
            snk.Write(keypath, s)
        }
    }
}

func RemoveAllSinks() {
    for _, snk := range sinks {
        if snk.needsClose {
            snk.writer.Close()
        }
    }
    sinks = make([]*sink, 0)
}

func (self *sink) String() string {
    return fmt.Sprintf("%v", self.keypathRegexp.String())
}

func (self *sink) MatchesKeypath(keypath string) bool {
    return self.keypathRegexp.MatchString(keypath)
}

func (self *sink) Write(prefix string, s string) {
    buf := new(bytes.Buffer)
    if self.writesPrefix {
        buf.WriteString(prefix)
        buf.WriteString(": ")
    }
    buf.WriteString(s)
    if len(s) > 0 && s[len(s)-1] != '\n' {
        buf.WriteByte('\n')
    }
    self.mu.Lock()
    defer self.mu.Unlock()
    self.writer.Write(buf.Bytes())
}

func (self *sink) SetWritesPrefix(pfx bool) {
    self.writesPrefix = pfx
}
