package main 

import (
	"fmt"
	"net"
	"flag"
	"encoding/binary"
	"os"
	"io"
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
		go sendFile(conn)
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

	defer conn.Close()

	message := "Hello"
	mSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(mSlice[0:], uint32(len(message)))
	_, err = conn.Write(mSlice)
	if err != nil {
		fmt.Println("Writing error")
		panic(err)
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Writing error")
		panic(err)
	}
	fileSizeSlice := make([]byte, 8)
	_, err = conn.Read(fileSizeSlice)
	if err != nil {
		fmt.Println("File size reading error")
		panic(err)
	}
	fileSize := binary.BigEndian.Uint64(fileSizeSlice[0:])
	newFile, err := os.Create("deliriumTWO.txt")
	if err != nil {
		fmt.Println ("File creating error")
		panic(err)
	}
	defer newFile.Close()

	receivedFile := make([]byte, fileSize)
	var receivedBytes uint64
	buffer := make([]byte, 1024)
	for {
		if receivedBytes >= fileSize {
			break
		}
		n, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		_ = copy(receivedFile[receivedBytes:], buffer)
		receivedBytes += uint64(n)
	}
	_, err = newFile.Write(receivedFile)
	if err != nil {

	}
}

func processConnection(conn *net.TCPConn) {
	msg := make([]byte, 4)
	_, err := conn.Read(msg)
	if err != nil {
		fmt.Println("Error reading ", err.Error())
	}
	n := binary.BigEndian.Uint32(msg[0:])
	decMsg := make([]byte, n)
	_, err = conn.Read(decMsg)
	if err != nil {
		fmt.Println("Error reading ", err.Error())
	}
	fmt.Printf("Message = %s\n", decMsg)
}

func sendFile(conn *net.TCPConn) {
	file, err := os.Create("deliriumONE.txt")
	if err != nil{
		fmt.Println("File error ", err.Error())
		panic(err)
	}
	defer file.Close()
	b := []byte("Blah-Blah")
	file.Write(b)
	file.Seek(0,0)
	err = os.Truncate("deliriumONE.txt", 8192)
	if err != nil {
		fmt.Println("Truncate error ", err.Error())
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
	}
	fileSizeSlice := make([]byte, 8)
	binary.BigEndian.PutUint64(fileSizeSlice[0:], uint64(fileInfo.Size()))
	_, err = conn.Write(fileSizeSlice)
	if err != nil {
		fmt.Println("File writing error")
		panic(err)
	}
	buffer := make([]byte, 1024)
	var count int
	for {
		if int64(count) >= fileInfo.Size(){
			break
		}
		nReadFromFile, err := file.Read(buffer)
		if err == io.EOF{
			break
		}
		nWritten, err := conn.Write(buffer[:nReadFromFile])
		if err != nil {
			panic(err)
		}
		count += nWritten
	}
}

 


















