package main 

import (
	"fmt"
	"net"
)

func main() {
	tcpAddress, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("tcpAddress error")
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		fmt.Println("Listener error")
		panic(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP error")
			panic(err)
		}
		go processConnection(conn)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		fmt.Println("DialTCP error")
		panic(err)
	}
	_, err = conn.Write([]byte ("Hello"))
	if err != nil {
		fmt.Println("Writing error")
		panic(err)
	}
}

func processConnection(conn *net.TCPConn){

	var msg []byte
	readmsg, err := conn.Read(msg)
	if err != nil {
		fmt.Println("Error reading ", err.Error())
	}
	fmt.Println("Recieved message is ", readmsg)
}