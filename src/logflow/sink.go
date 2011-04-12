package logflow

import (
    "os"
    "regexp"
    "strings"
)

type Sink interface {
    
}

type sink struct {
    keypathRegexp *regexp.Regexp
}

func NewSink(matches ...string) (*sink, os.Error) {
    re, err := regexp.Compile(strings.Join(matches, "|"))
    if err == nil {
        return &sink{
            keypathRegexp: re,
        }, nil
    }
    return nil, err
}

func (self *sink) String() string {
    return self.keypathRegexp.String()
}
