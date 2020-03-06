package main

import (
	"fmt"
	"os"

	"github.com/hpcloud/tail"
)

var config = tail.Config{
	ReOpen:    true,
	Follow:    true,
	MustExist: false,
	Poll:      true,
}

func main() {
	c := make(chan string)

	for _, f := range os.Args[1:] {
		go tailFile(f, c)
	}

	for {
		if <-c == "alarm" {
			fmt.Println("Alarm has been triggered on word \"ALARM\"")
		}
		fmt.Println(<-c)

	}

}

func tailFile(s string, c chan string) {
	t, err := tail.TailFile(s, config)
	if err != nil {
		panic(err)
	}
	for line := range t.Lines {
		c <- line.Text
	}

}
