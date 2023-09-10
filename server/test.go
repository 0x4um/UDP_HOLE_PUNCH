package main 

import (
	"fmt"
	"bufio"
	"os"
)

func main(){
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n\n\n\n\n\n>")
		value, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println(value[:len(value) - 1])
	}
}