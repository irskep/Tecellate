package main

import (
	"fmt"
	"os"
	"../../easynet"
)

func main() {
	fmt.Printf("testbot launched on %s\n", os.Args[0])
	listener := easynet.HostWithAddress(os.Args[0])
	defer listener.Close()
	conn := easynet.Accept(listener)
	defer conn.Close()
	
	conn.Write([]uint8("Response from " + os.Args[0]))
}
