package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

func getProblemsAnswers(csv [][]string) ([]string, []string) {
	var problems = []string{}
	var answers = []string{}
	for _, row := range csv {
		problems = append(problems, row[0])
		answers = append(answers, row[1])
	}
	return problems, answers
}
func readCsv(filePath string) ([]string, []string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	problems, answers := getProblemsAnswers(records)
	return problems, answers
}

func doQuiz(problems []string, answers []string, counterChannel chan int, quizEnd chan bool) {
	scanner := bufio.NewScanner(os.Stdin)
	counter := 0

	for i, prob := range problems {
		fmt.Println("What is: ", prob)
		userAns := ""
		for scanner.Scan() {
			userAns = scanner.Text()
			break
		}
		if userAns == answers[i] {
			counter += 1
			counterChannel <- counter
		}
	}
	score := fmt.Sprint(counter, "/", len(problems))
	fmt.Println("Your score is: ", score)
	quizEnd <- true
}

func main() {
	problems, answers := readCsv("problems.csv")
	counterChannel := make(chan int)
	counter := 0
	timeEnd := time.After(2 * time.Second)
	quizEnd := make(chan bool)
	go doQuiz(problems, answers, counterChannel, quizEnd)
	for {
		select {
		case <-counterChannel:
			counter += 1
		case <-timeEnd:
			score := fmt.Sprint(counter, "/", len(problems))
			fmt.Println("Your score is: ", score)
			return
		case <-quizEnd:
			return
		}
	}
}
