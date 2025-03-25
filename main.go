package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	CreateFile()
	inputChan := make(chan string)

	// Goroutine to read user input
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inputChan <- strings.TrimSpace(scanner.Text())
		}
	}()

	fmt.Println("TO DO Task Manager CLI (type 'exit' to quit)")
	fmt.Println("Usage: task [command] [options]")

	for {
		select {
		case input := <-inputChan:
			if input == "exit" {
				fmt.Println("Exiting Task Manager...")
				return
			}
			processCommand(input)
		case <-time.After(120 * time.Second): // Timeout message
			fmt.Println("Waiting for input...")
		}
	}
}

func processCommand(input string) {
	args := strings.Fields(input) // Split input by spaces

	if len(args) < 2 || args[0] != "task" {
		fmt.Println("Invalid command. Use 'task help' for usage.")
		return
	}
	command := args[1]
	options := args[2:]

	switch command {
	case "add":
		HandleAdd(options)
	case "delete":
		HandleDelete(options)
	case "list":
		HandleList(options)
	case "complete":
		HandleComplete(options)
	case "help":
		HandleHelp()
	default:
		fmt.Println("Unknown command. Use 'task help' for available commands.")
	}

}
