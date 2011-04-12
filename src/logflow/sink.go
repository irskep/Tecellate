package logflow

import (
    "io"
    "os"
    "regexp"
    "strings"
)

type Sink interface {
    
}

type sink struct {
    keypathRegexp *regexp.Regexp
    writer io.Writer
}

func NewSink(writer io.Writer, matches ...string) (*sink, os.Error) {
    re, err := regexp.Compile(strings.Join(matches, "|"))
    if err == nil {
        return &sink{
            keypathRegexp: re,
            writer: writer,
        }, nil
    }
    return nil, err
}

func (self *sink) String() string {
    return self.keypathRegexp.String()
}
