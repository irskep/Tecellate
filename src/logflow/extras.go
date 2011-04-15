package logflow

import (
    "bufio"
    "io"
    "os"
    "testing"
)

// TestWriter

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

// BufWriter

type BufWriter struct {
    f io.WriteCloser
    wr *bufio.Writer
}

func NewBufWriter(f io.WriteCloser) *BufWriter {
    return &BufWriter{f, bufio.NewWriter(f)}
}

func (self *BufWriter) Write(p []byte) (int, os.Error) {
    return self.wr.Write(p)
}

func (self *BufWriter) Close() os.Error {
    self.wr.Flush()
    return self.f.Close()
}
