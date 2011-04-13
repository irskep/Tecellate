package logflow

import (
    "bytes"
    "testing"
)

func TestHookup(t *testing.T) {
    w1 := new(bytes.Buffer)
    w2 := new(bytes.Buffer)
    w3 := new(bytes.Buffer)
    snk1, _ := NewSink(w1, "test1(/.*)?")
    snk2, _ := NewSink(w2, "test2(/.*)?")
    snk3, _ := NewSink(w3, "test1(/.*)?", "test2(/.*)?")
    src1 := NewSource("test1")
    src2 := NewSource("test2")
    
    var shouldBe string
    
    src1.Println("ABC")
    
    shouldBe = "test1/info: ABC\n"
    if !bytes.Equal([]byte(shouldBe), w1.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk1, shouldBe, w1.String())
    }
    shouldBe = ""
    if !bytes.Equal([]byte(shouldBe), w2.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk2, shouldBe, w2.String())
    }
    shouldBe = "test1/info: ABC\n"
    if !bytes.Equal([]byte(shouldBe), w3.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk3, shouldBe, w3.String())
    }
    
    src2.Println("DEF")
    shouldBe = "test1/info: ABC\n"
    if !bytes.Equal([]byte(shouldBe), w1.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk1, shouldBe, w1.String())
    }
    shouldBe = "test2/info: DEF\n"
    if !bytes.Equal([]byte(shouldBe), w2.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk2, shouldBe, w2.String())
    }
    shouldBe = "test1/info: ABC\ntest2/info: DEF\n"
    if !bytes.Equal([]byte(shouldBe), w3.Bytes()) {
        t.Errorf("%v mismatch:\n%v%v)", snk3, shouldBe, w3.String())
    }
}