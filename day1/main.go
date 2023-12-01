package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("expected to open file without errors. err:", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	// constant lookup
	numbersSet := map[rune]bool{
		'1': true,
		'2': true,
		'3': true,
		'4': true,
		'5': true,
		'6': true,
		'7': true,
		'8': true,
		'9': true,
		'0': true,
	}

	calibrationValues := make([]int64, 0)
	for scanner.Scan() {
		line := scanner.Text()
		queue := make([]rune, 0)

		for _, char := range line {
			if _, ok := numbersSet[char]; ok {
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
