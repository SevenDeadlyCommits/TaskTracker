package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/akamensky/argparse"
)

// Tracker The Tracker Struct
type Tracker struct {
	Name      string
	StartTime time.Time
}

func getDuration(t Tracker) time.Duration {
	return time.Since(t.StartTime)
}

func newTracker(name string) *Tracker {
	t := new(Tracker)
	t.Name = name
	t.StartTime = time.Now()
	return t
}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

func main() {

	// consoleReader := bufio.NewReader(os.Stdin)
	var taskTracker = new(Tracker)

	// Create new parser object
	parser := argparse.NewParser("print", "Prints provided string to stdout")
	// Create operation flag
	var taskName *string = parser.String("t", "track", &argparse.Options{Required: false, Help: "Start tracking a new task"})
	var statusUpdate *bool = parser.Flag("s", "status", &argparse.Options{Required: false, Help: "Show current tracked task status"})
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	if *statusUpdate {
		err = readGob("./task.gob", taskTracker)
		if err != nil {
			fmt.Println(err)
		} else {
			if taskTracker != nil {
				fmt.Printf("Currently tracking \"%s\": current duration: %s \n", taskTracker.Name, getDuration(*taskTracker).String())
			} else {
				fmt.Print("Not currently tracking a task")
			}
		}

	} else if len(*taskName) > 0 {
		taskTracker = newTracker(*taskName)
		err := writeGob("./task.gob", taskTracker)
		if err != nil {
			fmt.Println(err)
		}
	}
}
