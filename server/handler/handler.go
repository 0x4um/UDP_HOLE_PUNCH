package handler

import (
	"fmt"
)

func Exec(data string) bool {
	fmt.Println("hello from exec", data)
	return true
}