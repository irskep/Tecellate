package logflow

import (
    "bytes"
    "testing"
)

func checkLog(t *testing.T, snk Sink, target string, actual *bytes.Buffer) {
    if !bytes.Equal([]byte(target), actual.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk, target, actual.String())
    }
}

func TestHookup(t *testing.T) {
    w1 := new(bytes.Buffer)
    w2 := new(bytes.Buffer)
    w3 := new(bytes.Buffer)
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