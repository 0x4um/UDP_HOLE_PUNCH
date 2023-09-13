package handler

import (
	"fmt"
	"strings"
	"net"
)

func dialCheck(conn net.Conn){
	fmt.Println("checking for dial")
}

func function1(conn net.Conn) {
	conn.Write([]byte("return hello"))
	fmt.Println("function1 hello from function1")
}

func function2() {
	fmt.Println("function2 hello")
}

func Exec(data string, n int, conn net.Conn) {
	fmt.Println(data)
	fmt.Println(n)
	// functions := map[string]func(){
	// 	"function1": function1,
	// 	"function2": function2,
	// }
	// buffer := make([]byte, 1024)
	other := strings.TrimSpace(data)
	fmt.Println(other, " other other")
	switch other {
	case "findpeer":
		function1(conn)
	default:
		fmt.Println("def")
	}
}