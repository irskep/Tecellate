package main

import (
	"os"
	"../../easynet"
)

func main() {
	conn := easynet.HostWithAddress(os.Args[0])
	defer conn.Close()
	
	conn.Write([]uint8("ok"))
}
