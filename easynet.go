package easynet

import (
	"fmt"
	"os"
	"net"
	"log"
	"strings"
	"time"
)

func DieIfError(err os.Error, msg string) {
	if err != nil { log.Exit("", msg, " in coordinator: ", err) }
}

func HostWithAddress(addrString string) *net.TCPConn {
	fmt.Printf("Listening with address %s\n", addrString)
	addr, err := net.ResolveTCPAddr(addrString);
	DieIfError(err, "TCP address resolution error")
	listener, err := net.ListenTCP("tcp", addr);
	DieIfError(err, "Listening error")
	
	conn, err := listener.AcceptTCP();
	DieIfError(err, "TCP accept error")
	conn.SetKeepAlive(true)
	conn.SetReadTimeout(30000)
	
	return conn
}

func Dial(addrString string) *net.TCPConn {
	addr, err := net.ResolveTCPAddr(addrString);
	DieIfError(err, "TCP address resolution error")
	conn, err := net.DialTCP("tcp", nil, addr)
	for err != nil && strings.HasSuffix(err.String(), "connection refused") {
		time.Sleep(10000)
		conn, err = net.DialTCP("tcp", nil, addr)
	}
	DieIfError(err, "Dial error")
	return conn
}

func ReceiveFrom(conn *net.TCPConn) []byte {
	rcvd := make([]byte, 4096)
	size, err := conn.Read(rcvd)
	for err != nil && strings.HasSuffix(err.String(), "temporarily unavailable") {
		time.Sleep(10000)
		size, err = conn.Read(rcvd)
	}
	DieIfError(err, "Receive error")
	return rcvd[0:size]
}
