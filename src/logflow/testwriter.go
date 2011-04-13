package logflow

import (
    "os"
    "testing"
)

type TestWriter struct {
    t *testing.T
}

func NewTestWriter(t *testing.T) *TestWriter {
    return &TestWriter{t}
}

func (self *TestWriter) Write(p []byte) (int, os.Error) {
    self.t.Log(string(p[:len(p)-1]))
    return len(p), nil
}

func (self *TestWriter) Close() os.Error {
    return nil
}
