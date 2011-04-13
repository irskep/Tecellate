package logflow

import (
    "io"
    "os"
    "regexp"
    "strings"
)

type Sink interface {
    MatchesKeypath(string) bool
    Write([]byte)
}

var sinks []Sink = make([]Sink, 0)

type sink struct {
    keypathRegexp *regexp.Regexp
    writer io.Writer
}

func NewSink(w io.Writer, matches ...string) (*sink, os.Error) {
    re, err := regexp.Compile(strings.Join([]string{"^", strings.Join(matches, "|"), "$"}, ""))
    var theSink *sink = nil
    if err == nil {
        theSink = &sink{keypathRegexp: re, writer: w}
        sinks = append(sinks, theSink)
    }
    return theSink, err
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
    sinks = make([]Sink, 0)
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
