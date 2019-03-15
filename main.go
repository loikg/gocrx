package main

import (
	"fmt"
	"log"
	"os"

	"github.com/loikg/gocrx/cmd"
)

func init() {
	f, err := os.Create("log.txt")
	if err != nil {
		fmt.Printf("Failed to create log file: %v", err)
	}
	log.SetOutput(f)
}

func main() {
	cmd.Execute()
}