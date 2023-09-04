package main 

import (
	"fmt"
	"net"
	"strings"
)
func main(){
	clients := []string{"START"}
	
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


		go handleConn(conn, clients)

	}
}

func addToArray(value string, clients []string){
	clients = append(clients, value)
	fmt.Println("added", len(clients))
}


func handleConn(conn net.Conn, clients []string){
	defer conn.Close()

	fmt.Println(clients)
	
	buffer := make([]byte, 1024)

	for {

		n, err := conn.Read(buffer)
		if err != nil {
			addToArray("test", clients)
			fmt.Println("not reading buffer")
			return
		}
		data := string(buffer[:n])
		dataToCheck := data[:5]
		if strings.TrimSpace(dataToCheck) == "hello" {
			fmt.Println("hello world")
			clients = append(clients, data[6:])
			fmt.Println(clients)
		} else {
			fmt.Println("string check")
		}
	}
}