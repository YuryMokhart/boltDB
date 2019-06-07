package main 

import (
	"fmt"
	"net"
	"flag"
	"encoding/binary"
)

func main() {
	flagMode := flag.String("mode", "server", "start in server or client mode")
	flag.Parse()
	if *flagMode == "server"{
		server()
	} else if *flagMode == "client"{
		client()
	}
}

func server() {
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
	defer listener.Close()
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP error")
			panic(err)
		}
		go processConnection(conn)
	}
}

func client() {
	tcpAddress, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	if err != nil {
		fmt.Println("tcpAddress error")
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		fmt.Println("DialTCP error")
		panic(err)
	}
	message := "Hello"
	//fmt.Println(len(message))
	mSlice := make([]byte, len(message))
	binary.BigEndian.PutUint32(mSlice[0:], uint32(len(message)))
	fmt.Printf("%x\n", mSlice)
	_, err = conn.Write(mSlice)
	if err != nil {
		fmt.Println("Writing error")
		panic(err)
	}
}

func processConnection(conn *net.TCPConn){

	msg := make([]byte, 10)
	_, err := conn.Read(msg)
	if err != nil {
		fmt.Println("Error reading ", err.Error())
	}
	fmt.Println("Recieved message is ", msg)
	n := binary.BigEndian.Uint32(msg[0:])
	decodedMsg := make([]byte, 10)
	for i := 0; i <= int(n); i++{ // int(n) to be changed
		_, err := conn.Read(decodedMsg)
		if err != nil {
			fmt.Println("Error reading ", err.Error())
		}
	}
	//fmt.Println("Recieved message length is ", readMsgLength)
	//fmt.Println("Recieved message is ", msg)
	fmt.Println("Decoded message lenght = ", n)
	fmt.Println("Decoded message  = ", decodedMsg)
}






















