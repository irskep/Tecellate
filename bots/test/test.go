package main

import (
	"fmt"
	"os"
	"net"
	"log"
	"time"
)

func dieIfError(err os.Error, msg string) {
	if err != nil { log.Exit("", msg, " in bot: ", err) }
}

func setupConnection() *net.TCPConn {
	fmt.Printf("Bot address: %s\n", os.Args[0])
	addr, err := net.ResolveTCPAddr(os.Args[0]);
	dieIfError(err, "TCP address resolution error")
	listener, err := net.ListenTCP("tcp", addr);
	dieIfError(err, "Listening error")
	
	conn, err := listener.AcceptTCP();
	dieIfError(err, "TCP accept error")
	conn.SetKeepAlive(true)
	conn.SetReadTimeout(30000)
	
	time.Sleep(1000000)
	
	fmt.Printf("Bot waiting for connections on %s...\n", addr)
	
	return conn
}

func receive_from(conn *net.TCPConn) []byte {
	rcvd := make([]byte, 4096)
	size, err := conn.Read(rcvd)
	dieIfError(err, "Receive error")
	return rcvd[0:size]
}

func main() {
	conn := setupConnection()
	defer conn.Close()
	
	conn.Write([]uint8("ok"))
}
