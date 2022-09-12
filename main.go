package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	filename := flag.String("f", "examples/problems.csv", "use -f for a csv filename")
	timeoutSeconds := flag.Int("t", 30, "time in seconds before the timeout")

	flag.Parse()
	println("Welcome to the Quiz-Game")
	println("Filename: ", *filename)

	// open file using fileName
	file, err := os.Open(*filename)

	if err != nil {
		exit(fmt.Sprintf("failed to open file:%s", *filename))
	}

	r := csv.NewReader(file)

	records, err := r.ReadAll()
	if err != nil {
		exit("Failed to read CSV file")
	}

	c := make(chan int)
	timer := time.NewTimer(time.Second * time.Duration(*timeoutSeconds))

	fmt.Println("timer fired with:", *timeoutSeconds, "seconds")
	go quiz(records, c, timer)

	answers := <-c

	fmt.Println()
	fmt.Println("Total number of correct answers questions:", answers, "/", len(records))
}

func quiz(records [][]string, c chan int, timer *time.Timer) {
	correctAnswers := 0

	for _, record := range records {

		answerCh := make(chan string)
		go func() {
			fmt.Println(record[0])
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Answer: ")

			text, err := reader.ReadString('\n')
			if err != nil {
				// this will exit the whole program, if a user enters a invalid input
				exit(fmt.Sprintf("Error reading from STdin: %v", err))
			}
			text = strings.TrimSuffix(text, "\n")
			answerCh <- text
		}()

		select {
		case <-timer.C:
			c <- correctAnswers
			return
		case answer := <-answerCh:
			if answer == record[1] {
				correctAnswers++
			}

		}
	}

	c <- correctAnswers
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
