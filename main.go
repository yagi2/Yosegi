package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("👋 Welcome to your CLI, yagi2!")
	args := os.Args[1:]
	if len(args) > 0 {
		fmt.Printf("You passed: %v\n", args)
	} else {
		fmt.Println("No args passed.")
	}
}
