package wifi

import "fmt"
import "hash/crc32"
import . "byteslice"

type TTL uint8
const DEFAULT_TTL = 128
const SEND_TTL = 10

//                                      TTL  CRC32  ADDRS
const MessageBodySize = PacketBodySize - 1  -  4  -   8

type Message struct {
    Message ByteSlice
    DestAddr uint32
    FromAddr uint32
    TTL TTL
    SendTTL TTL
    checksum ByteSlice
}

func NewMessage(msg ByteSlice, from, dest uint32) *Message {
    bytes := make(ByteSlice, MessageBodySize)
    copy(bytes, msg)
    return &Message{
        Message:bytes,
        DestAddr:dest,
        FromAddr:from,
        TTL:DEFAULT_TTL,
        SendTTL:SEND_TTL,
    }
}

func MakeMessage(msg ByteSlice) *Message {
    bytes := make(ByteSlice, PacketBodySize)
    copy(bytes, msg)
    dest_addr := bytes[0:4]
    from_addr := bytes[4:8]
    ttl       := bytes[8:9]
    message   := bytes[9:len(bytes)-4]
    crc       := bytes[len(bytes)-4:]
    return &Message{
        Message:message,
        DestAddr:dest_addr.Int32(),
        FromAddr:from_addr.Int32(),
        TTL:TTL(ttl.Int8()),
        SendTTL:SEND_TTL,
        checksum:crc,
    }
}

func (self *Message) Body() ByteSlice {
    return self.Message
}

func (self *Message) DecTTL() {
    if self.TTL > 0 { self.TTL -= 1 }
    if self.SendTTL > 0 { self.SendTTL -= 1 }
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
    bytes     := make(ByteSlice, PacketBodySize)
    dest_addr := bytes[0:4]
    from_addr := bytes[4:8]
    ttl       := bytes[8:9]
    msg       := bytes[9:len(bytes)-4]
    copy(dest_addr, ByteSlice32(self.DestAddr))
    copy(from_addr, ByteSlice32(self.FromAddr))
    copy(ttl, ByteSlice8(uint8(self.TTL)))
    copy(msg, self.Message)
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
    return fmt.Sprintf("<Message from:%v to:%v ttl:%v>", self.FromAddr, self.DestAddr, self.TTL)
}
