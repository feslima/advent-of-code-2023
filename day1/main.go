package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// constant lookup
var digitSet map[rune]bool = map[rune]bool{
	'1': true,
	'2': true,
	'3': true,
	'4': true,
	'5': true,
	'6': true,
	'7': true,
	'8': true,
	'9': true,
}

/*
There are instances of overlapping numbers (e.g: twone that must be
interpreted as 21, etc). Therefore, this maps the overlaps by linking
numbers that start/ends with same letters. Then, replace the matches in
an way that breaks the overlap.
*/
var overlappingNumbers map[string]string = map[string]string{
	"one":   "o1e", // matches with eight
	"two":   "t2o", // matches with one
	"three": "t3e", // matches with eight
	"four":  "4",
	"five":  "5e", // matches with eight
	"six":   "6",
	"seven": "7n",  // matches with nine
	"eight": "e8t", // matches with two, three
	"nine":  "n9e", // matches with eight
}

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

	calibrationValues := make([]int64, 0)
	for scanner.Scan() {
		line := scanner.Text()

		for number, digit := range overlappingNumbers {
			if i := strings.Index(line, number); i >= 0 {
				line = strings.ReplaceAll(line, number, digit)
			}
		}
		queue := make([]rune, 0)

		for _, char := range line {
			if _, ok := digitSet[char]; ok {
				queue = append(queue, char)
			}
		}

		if len(queue) == 0 {
			// found nothing, proceed to next line
			continue
		} else if len(queue) == 1 {
			// single char must be treated as double char
			queue = append(queue, queue[0])
		}

		val, err := strconv.ParseInt(fmt.Sprintf("%c%c", queue[0], queue[len(queue)-1]), 10, 0)
		if err != nil {
			// for some reason, failed to parse string representation to integer
			continue
		}
		calibrationValues = append(calibrationValues, val)
	}

	sum := int64(0)
	for _, v := range calibrationValues {
		sum += v
	}

	fmt.Printf("Sum: %d\n", sum)
}
