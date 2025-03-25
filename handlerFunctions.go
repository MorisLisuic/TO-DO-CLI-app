package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

const Filename = "data.csv"

func CreateFile() {
	if _, err := os.Stat(Filename); os.IsNotExist(err) {
		file, err := os.Create(Filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		record := [][]string{
			{"Id", "Task", "Created", "Done"},
		}

		for _, value := range record {
			err := writer.Write(value)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func HandleAdd(options []string) {
	if len(options) > 0 {
		if checkHelpFlag(options[0]) {
			fmt.Println("Usage: task add [string you want to save]")
			return
		}
	}
	file, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var lastRecord []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lastRecord = record
	}
	var newRecord []string
	if lastRecord[0] == "Id" {
		newRecord = append(newRecord, "1")
	} else {
		newId, _ := strconv.ParseInt(lastRecord[0], 10, 32)
		newRecord = append(newRecord, strconv.FormatInt(newId+1, 10))
	}
	task := strings.Join(options, " ")
	newRecord = append(newRecord, task)
	currentTime := time.Now().Format("2006-01-02/15:04:05")
	newRecord = append(newRecord, currentTime)
	newRecord = append(newRecord, "False")

	file, err = os.OpenFile(Filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.Write(newRecord)
	writer.Flush()
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	fmt.Println("")
	fmt.Fprintln(w, "Id\tTask\tCreated\tDone")
	consoleString := ""
	for index, value := range newRecord {
		if index == 0 {
			consoleString += value
		} else {
			consoleString += "\t" + value
		}
	}
	fmt.Fprintln(w, consoleString)
	w.Flush()

}

func HandleDelete(options []string) {
	if len(options) > 0 {
		if checkHelpFlag(options[0]) {
			fmt.Println("Usage: 'task delete [Id of task to delete]'")
			return
		}
	}
	file, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var records [][]string
	reader := csv.NewReader(file)
	foundFlag := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if record[0] == options[0] {
			foundFlag = true
		} else {
			records = append(records, record)
		}

	}
	file, err = os.OpenFile(Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
	fmt.Println("")
	if foundFlag {
		fmt.Println("Task deleted")
	} else {
		fmt.Println("Task with Id ", options[0], " does not exist")
	}
}

func HandleList(options []string) {
	if len(options) > 0 {
		if checkHelpFlag(options[0]) {
			fmt.Println("Usage: 'task list' outputs a table with all tasks")
			return
		}
	}
	file, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	fmt.Println("")
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		consoleString := ""
		for index, value := range record {
			if index == 0 {
				consoleString += value
			} else {
				consoleString += "\t" + value
			}
		}
		fmt.Fprintln(w, consoleString)
	}
	w.Flush()
}

func HandleComplete(options []string) {
	if len(options) > 0 {
		if checkHelpFlag(options[0]) {
			fmt.Println("Usage: 'task complete [Id of the task to mark as complete]'")
			return
		}
	}
	file, err := os.Open(Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var records [][]string
	reader := csv.NewReader(file)
	foundFlag := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if record[0] == options[0] {
			record[3] = "True"
			foundFlag = true
		}
		records = append(records, record)
	}
	file, err = os.OpenFile(Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, record := range records {
		if err := writer.Write(record); err != nil {
			fmt.Println("Error writing record:", err)
			return
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}
	fmt.Println("")
	if foundFlag {
		fmt.Println("Task marked as complete")
	} else {
		fmt.Println("Task with Id ", options[0], " does not exist")
	}
}
func HandleHelp() {
	fmt.Println("Usage:")
	fmt.Println("	tasks [command] [flag](optional)")
	fmt.Println("")
	fmt.Println("Available commands:")

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	fmt.Fprintln(w, "	add\tAdd a new task to the todo list")
	fmt.Fprintln(w, "	complete\tSet a task as completed")
	fmt.Fprintln(w, "	delete\tDelete a task")
	fmt.Fprintln(w, "	list\tList all tasks")
	w.Flush()

	fmt.Println("")
	fmt.Println("Flags:")
	fmt.Fprintln(w, "	-h, --help\tShow this help")
	w.Flush()
	fmt.Println("")
	fmt.Println("Use 'task [command] --help' or 'Use 'task [command] --h' for information about a command")
}

func checkHelpFlag(flag string) bool {
	if flag == "-h" || flag == "--help" {
		return true
	}
	return false
}
