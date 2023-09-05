package main 

import (
	"fmt"
	"net"
	"strings"
)
func main(){
	
	ln, err := net.Listen("tcp", ":12000")
	if err != nil {
		fmt.Println("error starting")
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("accepting error")
			continue
		}


		go handleConn(conn)

	}
}
func handleConn(conn net.Conn){
	defer conn.Close()

	
	buffer := make([]byte, 1024)

	for {

		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("not reading buffer")
			return
		}
		data := string(buffer[:n])
		dataToCheck := data[:5]
		if strings.TrimSpace(dataToCheck) == "hello" {
			fmt.Println("hello world", data)
			localAddr := conn.LocalAddr().String()
			_, err := conn.Write([]byte(localAddr))
			if err != nil {
				fmt.Println("unable to write")
				return
			}
		} else {
			fmt.Println("string check")
		}
	}
}