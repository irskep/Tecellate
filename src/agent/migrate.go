package agent

import "fmt"
import . "byteslice"

type Migrate struct {
    address []byte
}

func NewMigrate(addr []byte) *Migrate {
    self := new(Migrate)
    self.address = addr
    return self
}


func MakeMigrate(bytes ByteSlice) *Migrate {
    return NewMigrate(bytes)
}

func (self *Migrate) String() string {
    return fmt.Sprintf("Migrate to %d", string(self.address))
}

func (self *Migrate) Bytes() ByteSlice {
    return self.address
}
