package main

import (
	"fmt"
	"encoding/csv"
	//"io"
	//"log"
	"time"
	"os"
	"bufio"
	"flag"
	"math/rand"
)	

func main() {
	problemsFile := flag.String("csv","problems.csv","CSV File in the format of 'question,answer'")
	limit := flag.Int("limit",30,"time limit for quiz to expire in seconds")
	shuffle := flag.Bool("shuffle",false,"Shuffle the order of Questions")
	flag.Parse()
	
	csvFile, fileError := os.Open(*problemsFile)
	if fileError != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *problemsFile))
	}
	r := csv.NewReader(csvFile)

	record, readerError := r.ReadAll()
	
	if *shuffle {
		rand.Shuffle(len(record), func(i,j int) {
			record[i], record[j] = record[j], record[i]
		})
	}
	
	if readerError != nil {
		exit("Failed to read the file which you have provided")
	}

	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	correctAnswer := 0
	totalQuestions := len(record)
	
	scanner := bufio.NewScanner(os.Stdin)
	problemLoop:
	for _, row := range record {
		fmt.Print(row[0]+": ")
		answerCh := make(chan string)
		go func() {
			scanner.Scan()
			answer := scanner.Text()
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			fmt.Println()
			//fmt.Printf("\nTimes Up, Scored %d out of %d \n", correctAnswer, totalQuestions)
		break problemLoop
		case answer := <-answerCh:
			if row[1] == answer {
				correctAnswer++;	
			}
		}
	}
	fmt.Println("Answered",correctAnswer,"out of",totalQuestions)
}

func exit(msg string) {
	fmt.Printf(msg)
	os.Exit(1)
}
