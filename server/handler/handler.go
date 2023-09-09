package handler

import (
	"fmt"
	"strings"
	"net"
)

func function1() {
	fmt.Println("function1 hello")
}

func function2() {
	fmt.Println("function2 hello")
}

func Exec(data string, n int, conn net.Conn) {
	fmt.Println(data)
	fmt.Println(n)
	functions := map[string]func(){
		"function1": function1,
		"function2": function2,
	}
	// buffer := make([]byte, 1024)
	other := strings.TrimSpace(data)
	fmt.Println(other, " other other")
	if fn, exists := functions[other]; exists {
		fn()
	} else {
		fmt.Println("func not found")
	}
}