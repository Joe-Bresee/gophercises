package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	csvFile := flag.String("csv", "problems.csv", "csv in format Q:A")
	timeLimit := flag.Int("limit", 30, "time limit in seconds")
	flag.Parse()

	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	r := csv.NewReader(file)
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	numRight := 0
	totalNum := 0

quizLoop:
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		if len(record) < 2 {
			continue
		}

		question := strings.TrimSpace(record[0])
		answer := strings.TrimSpace(record[1])

		answerCh := make(chan string)

		fmt.Printf("%s = ", question)
		go func() {
			var userAnswer string
			fmt.Scanln(&userAnswer)
			answerCh <- strings.TrimSpace(userAnswer)
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			break quizLoop
		case userAnswer := <-answerCh:
			if userAnswer == answer {
				numRight++
			}
			totalNum++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", numRight, totalNum)
}
