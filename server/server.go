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
	var clients [3]string
	for k := 0; k < len(clients); k++{
		clients[k] = "empty"
	}
	buffer := make([]byte, 1024)

	for {

		n, err := conn.Read(buffer)
		if err != nil {
			for j := 0; j < len(clients); j++{
				fmt.Println(clients[j])
			}
			fmt.Println("not reading buffer")
			return
		}
		data := string(buffer[:n])
		dataToCheck := data[:5]
		if strings.TrimSpace(dataToCheck) == "hello" {
			fmt.Println("hello world")
			if clients[1] == "empty" {
				fmt.Println("0 empty")
				clients[1] = data[6:]
			} else {
				fmt.Println("0 not empty")
				clients[2] = data[6:]
			}
			fmt.Println(clients[0] + "whereami")
		} else {
			fmt.Println("string check")
		}
	}
}