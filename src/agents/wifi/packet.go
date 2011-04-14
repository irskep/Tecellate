package wifi

import "fmt"
import "hash/crc32"
import "coord/game"
import . "byteslice"

type Command uint8
var Commands map[string]Command
var cmdsr []string

type Packet struct {
    pkt [game.MessageLength]byte
}


// init functions --------------------------------------------------------------
func init() {
    Commands = make(map[string]Command)
    cmdsr = []string{
        "HELLO", "ACK", "NAK",
    }
    for i, cmd := range cmdsr {
        Commands[cmd] = Command(i)
    }
}


// packet methods --------------------------------------------------------------
func NewPacket(cmd Command) *Packet {
    self := new(Packet)
    self.pkt[0] = byte(cmd)
    return self
}

func MakePacket(bytes ByteSlice) *Packet {
    self := new(Packet)
    for i := 0; i < game.MessageLength; i++ {
        self.pkt[i] = bytes[i]
    }
    return self
}

func (self *Packet) Cmd() (bool, Command, string) {
    cmd := Command(self.pkt[0])
    if cmd < Command(len(cmdsr)) {
        return true, cmd, cmdsr[cmd]
    }
    return false, 0, ""
}

func (self *Packet) Bytes() ByteSlice {
    pkt := self.pkt[:]
    copy(pkt[len(pkt)-4:], self.ComputeChecksum())
    return pkt
}

func (self *Packet) ComputeChecksum() ByteSlice {
    return ByteSlice32(crc32.ChecksumIEEE(self.pkt[:len(self.pkt)-4]))
}

func (self *Packet) ValidateChecksum() bool {
    checksum := ByteSlice(self.pkt[len(self.pkt)-4:])
    return checksum.Eq(self.ComputeChecksum())
}

func (self *Packet) String() string {
    var command string
    if ok, _, name := self.Cmd(); ok {
        command = name
    } else {
        command = "Unknown"
    }
    return fmt.Sprintf("<Packet cmd:%v %v %v>", command, self.pkt[len(self.pkt)-4:], self.ValidateChecksum())
}
