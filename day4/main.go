package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Card map[int]struct {
	winners []string
	count   int
}

func main() {
	filename := os.Args[1]
	// filename := "day4-example.txt"

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("expected to open file without errors. err:", err)
		os.Exit(1)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	cardRegExp := regexp.MustCompile(`Card\s+(\d+):\s+([\d\s]+)\s*\|\s*([\d\s]+)`)
	cards := make(Card)
	for scanner.Scan() {
		line := scanner.Text()

		cardCapture := cardRegExp.FindStringSubmatch(line)
		if cardCapture == nil || len(cardCapture) < 4 {
			continue
		}

		cardNum, err := strconv.ParseInt(string(cardCapture[1]), 10, 0)
		if err != nil {
			continue
		}

		winningNumbers := strings.Split(cardCapture[2], " ")
		cardNumbers := strings.Split(cardCapture[3], " ")

		numbers := map[string]bool{}
		for _, n := range winningNumbers {
			_, err := strconv.ParseInt(n, 10, 0) // only add parsable numbers
			if err != nil {
				continue
			}
			numbers[n] = true
		}

		found := []string{}
		for _, n := range cardNumbers {
			if _, exists := numbers[n]; exists {
				found = append(found, n)
			}
		}

		cards[int(cardNum)] = struct {
			winners []string
			count   int
		}{winners: found, count: 1}
	}

	points := 0
	/* we have to iterate via conventional loop instead of
	map iteration because the order of indexing matters and
	golang map doesn't keep order of keys.
	*/
	for i := 1; i < len(cards); i++ {
		currentCard, ok := cards[i]
		if !ok {
			continue
		}

		n := len(currentCard.winners)
		if n > 0 {
			points += int(math.Pow(2, float64(n-1)))
		}

		for j := i + 1; j <= i+n; j++ {
			nextCard, ok := cards[j]
			if ok {
				cards[j] = struct {
					winners []string
					count   int
				}{winners: nextCard.winners, count: nextCard.count + currentCard.count}
			}
		}

	}
	fmt.Printf("Sum: %d\n", points)

	total := 0
	for _, e := range cards {
		total += e.count
	}
	fmt.Printf("total: %d\n", total)
}
