package logflow

import (
    "bytes"
    "os"
    "testing"
)

/*
    Helpers
*/

type TestWriter struct {
    Contents []byte
}

func NewTestWriter() *TestWriter {
    return &TestWriter{
        Contents: make([]byte, 0),
    }
}

func (self *TestWriter) Write(p []byte) (n int, err os.Error) {
    self.Contents = bytes.Join([][]byte{self.Contents, p}, []byte{})
    return len(p), nil
}

func (self *TestWriter) String() string {
    return string(self.Contents)
}

/*
    Tests
*/

func TestTestWriter(t *testing.T) {
    w := NewTestWriter()
    testString := []byte("ABC\n")
    w.Write(testString)
    if !bytes.Equal(testString, w.Contents) {
        t.Fatalf("Strings don't match! (%v, %v)", testString, w.Contents)
    }
    testString2 := []byte("DEF\n")
    w.Write(testString2)
    concattedTestString := []byte("ABC\nDEF\n")
    if !bytes.Equal(concattedTestString, w.Contents) {
        t.Fatalf("Strings don't match! (%v, %v)", testString, w.Contents)
    }
}

func TestSourceInstantiate(t *testing.T) {
    src := NewSource("test.info")
    t.Log(src)
}

func TestSinkInstantiate(t *testing.T) {
    sink, err := NewSink(nil, "test/info", "test/debug")
    t.Log(sink, err)
}

func TestHookup(t *testing.T) {
    w := NewTestWriter()
    sink, err := NewSink(w, "test/info")
    t.Log(sink, err)
    src := NewSource("test/info")
    src.Println("Hello!")
}