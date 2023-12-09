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

type Card map[int][]string

func main() {
	filename := os.Args[1]

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

		cards[int(cardNum)] = found
	}

	points := 0
	for _, nums := range cards {
		n := len(nums)
		if n > 0 {
			points += int(math.Pow(2, float64(n-1)))
		}
	}
	fmt.Printf("Sum: %d\n", points)
}
