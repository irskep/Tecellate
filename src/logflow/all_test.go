package logflow

import (
    "bytes"
    "os"
    "testing"
)

type ClosableBuffer struct {
    buf *bytes.Buffer
}

func NewClosableBuffer() *ClosableBuffer {
    return &ClosableBuffer{new(bytes.Buffer)}
}

func (self *ClosableBuffer) Write(p []byte) (int, os.Error) {
    return self.buf.Write(p)
}

func (self *ClosableBuffer) Close() os.Error {
    return nil
}

func (self *ClosableBuffer) String() string {
    return self.buf.String()
}

func (self *ClosableBuffer) Bytes() []byte {
    return self.buf.Bytes()
}

func checkLog(t *testing.T, snk Sink, target string, actual *ClosableBuffer) {
    if !bytes.Equal([]byte(target), actual.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk, target, actual.String())
    }
}

func TestHookup(t *testing.T) {
    w1 := NewClosableBuffer()
    w2 := NewClosableBuffer()
    w3 := NewClosableBuffer()
    snk1, _ := NewSink(w1, "test1/.*")
    snk2, _ := NewSink(w2, "test2/.*")
    snk3, _ := NewSink(w3, "test1/.*", "test2/.*")
    src1 := NewSource("test1")
    src2 := NewSource("test2")
    
    src1.Println("ABC")
    checkLog(t, snk1, "test1/info: ABC\n", w1)
    checkLog(t, snk2, "", w2)
    checkLog(t, snk3, "test1/info: ABC\n", w3)
    
    src2.Println("DEF")
    checkLog(t, snk1, "test1/info: ABC\n", w1)
    checkLog(t, snk2, "test2/info: DEF\n", w2)
    checkLog(t, snk3, "test1/info: ABC\ntest2/info: DEF\n", w3)
}

func TestTestWriter(t *testing.T) {
    NewSink(NewTestWriter(t), ".*")
    src1 := NewSource("testwriterTest")
    src1.Println("This test fails on purpose.")
    // t.Fatal("On purpose, I say!")
}
