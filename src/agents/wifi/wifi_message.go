package wifi

import "fmt"
import . "byteslice"

type TTL uint8
const DEFAULT_TTL = 128
const SEND_TTL = 10

//                                      TTL  CRC32  ADDRS
const MessageBodySize = PacketBodySize - 1  -  4  -   8

type Message struct {
    Message ByteSlice
    TTL TTL
    DestAddr uint32
    FromAddr uint32
    checksum ByteSlice
    sendTTL TTL
}

func NewMessage(msg ByteSlice, from, dest uint32) *Message {
    bytes := make(ByteSlice, MessageBodySize)
    copy(bytes, msg)
    return &Message{
        Message:bytes,
        TTL:DEFAULT_TTL,
        sendTTL:SEND_TTL,
        DestAddr:dest,
        FromAddr:from,
    }
}

func (self *Packet) ComputeChecksum() ByteSlice {
    bytes     := self.Bytes()
    crc       := bytes[len(bytes)-4:]
    return crc
}

func (self *Packet) ValidateChecksum() bool {
    bytes     := self.ComputeChecksum()
    crc       := bytes[len(bytes)-4:]
    return self.checksum.Eq(self.ComputeChecksum())
}

func (self *Message) body_bytes() ByteSlice {
    bytes     := make(ByteSlice, PacketBodySize)
    body      := bytes[:len(bytes)-4]
    dest_addr := bytes[0:4]
    from_addr := bytes[4:8]
    ttl       := bytes[8:9]
    msg       := bytes[9:len(bytes)-4]
    copy(dest_addr, ByteSlice32(self.DestAddr))
    copy(from_addr, ByteSlice32(self.FromAddr))
    copy(ttl, ByteSlice8(self.TTL))
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
    return fmt.Sprintf("<Message from:%v to:%v ttl:%v>", self.FromAddr, self.DestAddr, self.TTL)
}
