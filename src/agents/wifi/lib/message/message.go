package message

import "fmt"
import "hash/crc32"
import . "byteslice"

import . "agents/wifi/lib/datagram"

type SEQUENCE uint16

//                                        SQC   ACK   crc32
const MessageBodySize = DataGramBodySize - 2  -  2  -   4

type Message struct {
    Message ByteSlice
    Sequence SEQUENCE
    Acknowledge SEQUENCE
    checksum ByteSlice
}

func NewMessage(msg ByteSlice, seq, ack SEQUENCE) *Message {
    bytes := make(ByteSlice, MessageBodySize)
    copy(bytes, msg)
    return &Message{
        Message:msg,
        Sequence:seq,
        Acknowledge:ack,
    }
}

func MakeMessage(msg ByteSlice) *Message {
    bytes := make(ByteSlice, DataGramBodySize)
    copy(bytes, msg)
    seq       := ByteSlice(bytes[0:2])
    ack       := ByteSlice(bytes[2:4])
    message   := bytes[4:len(bytes)-4]
    crc       := bytes[len(bytes)-4:]
    return &Message{
        Message:message,
        Sequence:SEQUENCE(seq.Int16()),
        Acknowledge:SEQUENCE(ack.Int16()),
        checksum:crc,
    }
}

func (self *Message) IsAck() bool {
    return self.Sequence + 1 == self.Acknowledge
}

func (self *Message) Body() ByteSlice {
    return self.Message
}

func (self *Message) ComputeChecksum() ByteSlice {
    bytes     := self.Bytes()
    crc       := bytes[len(bytes)-4:]
    return crc
}

func (self *Message) ValidateChecksum() bool {
    return self.checksum.Eq(self.ComputeChecksum())
}

func (self *Message) body_bytes() ByteSlice {
    bytes     := make(ByteSlice, DataGramBodySize)
    seq       := bytes[0:2]
    ack       := bytes[2:4]
    message   := bytes[4:len(bytes)-4]
    copy(seq, ByteSlice16(uint16(self.Sequence)))
    copy(ack, ByteSlice16(uint16(self.Acknowledge)))
    copy(message, self.Message)
    return bytes
}

func (self *Message) Bytes() ByteSlice {
    bytes     := self.body_bytes()
    body      := bytes[:len(bytes)-4]
    crc       := bytes[len(bytes)-4:]
    copy(crc, ByteSlice32(crc32.ChecksumIEEE(body)))
    return bytes
}

func (self *Message) String() string {
    if self == nil { return "<nil Message>" }
    return fmt.Sprintf("<Message seq:%v ack:%v>", self.Sequence, self.Acknowledge)
}
