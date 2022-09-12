package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	defaultFileName = "examples/problems.csv"
	defaultTimeout  = 30
	defaultShuffle  = false
)

func main() {

	filename := flag.String("f", defaultFileName, "use to specify a CSV filename")
	timeoutSeconds := flag.Int("t", defaultTimeout, "time in seconds for the quiz timeout")
	shuffle := flag.Bool("shuffle", defaultShuffle, " boolean to enable shuffling of the questions")

	flag.Parse()

	println("Welcome to the Quiz-Game")

	// open file using fileName
	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("failed to open file:%s", *filename))
	}
	// create CSV reader
	r := csv.NewReader(file)
	// read all records in CSV and return as slice[][]string
	records, err := r.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to read from CSV file: %v", err))
	}

	// bonus, shuffle mode to randomize question order
	if *shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(records), func(i, j int) {
			records[i], records[j] = records[j], records[i]
		})
	}

	c := make(chan int)
	timer := time.NewTimer(time.Second * time.Duration(*timeoutSeconds))

	fmt.Println("Timer started with:", *timeoutSeconds, "seconds of timeout. Hurry up!")
	go quiz(records, c, timer)

	answers := <-c

	fmt.Printf("Total number of correct answers: %d/%d\n", answers, len(records))
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
			text = strings.TrimSpace(text)
			text = strings.ToLower(text)

			answerCh <- text
		}()

		select {
		case <-timer.C:
			c <- correctAnswers
			return
		case answer := <-answerCh:
			if answer == strings.ToLower(record[1]) {
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
