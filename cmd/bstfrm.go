package main

import (
	"bstfrm"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to bstfrm.")

	for {
		fmt.Print("# ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		ast, err := bstfrm.Parse(text)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		m := bstfrm.NewMachine(ast)
		m.Run()
		fmt.Println()
		fmt.Println("ok")
	}
}
