package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type quiz struct {
	question string
	answer   string
}

func ticking(countDownSecond int, c chan bool) {
	time.Sleep(time.Second * time.Duration(countDownSecond))
	c <- true
}

func read(questionBanks []*quiz, c chan bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < len(questionBanks); i++ {
		fmt.Print(questionBanks[i].question, "=")
		scanner.Scan()
		res := scanner.Text()
		if res == questionBanks[i].answer {
			c <- true
		}

		c <- false
	}

}

func main() {
	// flags
	fileLocation := flag.String("filepath", "problems.csv", "filepath for parsing")
	limit := flag.Int("limit", 5, "time limit")

	flag.Parse()

	// Read csv
	var questionBanks []*quiz

	file, err := os.Open(*fileLocation)
	defer file.Close()
	if err != nil {
		fmt.Println("File failed to read")
	}
	reader := bufio.NewReader(file)

	for {
		n, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		// Process
		arr := strings.Split(strings.Trim(string(n), "\n"), ",")

		questionBanks = append(questionBanks, &quiz{question: arr[0], answer: arr[1]})
	}

	// // count down elapsed
	done := make(chan bool)
	judgement := make(chan bool)

	go ticking(*limit, done)
	go read(questionBanks, judgement)

	score := 0

	// communicate if tmer timeout, return with messages
	for {
		select {
		case <-done:
			fmt.Println("\nYour time running out. Your score is", score)
			close(judgement)
			return
		case res := <-judgement:
			// evaluate corectness
			if res == true {
				score++
			}
		}
	}
}
