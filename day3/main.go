package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	'0': true,
}

func getLeftMostIndex(start int, line string) int {
	if start == 0 {
		return start
	}

	pointer := start
	for j := start; j >= 0; j-- {
		char := rune(line[j])
		pointer = j
		if _, ok := digitSet[char]; !ok {
			pointer++
			break
		}
	}
	return pointer
}

func getRightMostIndex(start int, line string) int {
	if start == (len(line) - 1) {
		return start
	}

	pointer := start
	for j := start; j < len(line); j++ {
		char := rune(line[j])
		pointer = j
		if _, ok := digitSet[char]; !ok {
			break
		}
	}
	return pointer
}

func getNeighborNumbers(cellRow int, cellCol int, lines []string, visits [][]bool) []string {
	height := len(lines)
	width := len(lines[0])
	var rowTopOffset int
	var rowBottomOffset int
	if cellRow == 0 {
		rowTopOffset = 0
		rowBottomOffset = 1
	} else if cellRow == (height - 1) {
		rowTopOffset = -1
		rowBottomOffset = 0
	} else {
		rowTopOffset = -1
		rowBottomOffset = 1
	}

	var colLeftOffset int
	var colRightOffset int
	if cellCol == 0 {
		colLeftOffset = 0
		colRightOffset = 1
	} else if cellCol == (width - 1) {

		colLeftOffset = -1
		colRightOffset = 0
	} else {
		colLeftOffset = -1
		colRightOffset = 1
	}

	indexes := [][]int{
		{cellRow + rowTopOffset, cellCol + colLeftOffset},
		{cellRow + rowTopOffset, cellCol},
		{cellRow + rowTopOffset, cellCol + colRightOffset},
		{cellRow, cellCol + colLeftOffset},
		{cellRow, cellCol + colRightOffset},
		{cellRow + rowBottomOffset, cellCol + colLeftOffset},
		{cellRow + rowBottomOffset, cellCol},
		{cellRow + rowBottomOffset, cellCol + colRightOffset},
	}
	numbers := make([]string, 0)
	for _, pos := range indexes {
		row := pos[0]
		col := pos[1]

		c := rune(lines[row][col])
		v := visits[row][col]
		if _, ok := digitSet[c]; ok && !v {
			// number found, sweep left and right
			left := getLeftMostIndex(col, lines[row])
			right := getRightMostIndex(col, lines[row])

			var number string
			if right == (width - 1) {
        // edge case
				number = lines[row][left:(right + 1)]
			} else {
				number = lines[row][left:right]
			}

			numbers = append(numbers, strings.ReplaceAll(number, ".", ""))

			for j := left; j < right; j++ {
				visits[row][j] = true
			}
		}
	}

	visits[cellRow][cellCol] = true
	return numbers
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

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	nLines := len(lines)
	width := len(lines[0])

	visits := make([][]bool, nLines)
	for i := 0; i < nLines; i++ {
		visits[i] = make([]bool, width)
		for j := 0; j < width; j++ {
			visits[i][j] = false
		}
	}

	collectedSymbols := make([]struct {
		x         int
		y         int
		c         string
		neighbors []string
	}, 0)
	collectedGears := make([]struct {
		numbers []string
	}, 0)
	for i := 0; i < nLines; i++ {
		line := lines[i]
		for j := 0; j < width; j++ {
			char := rune(line[j])
			if char != '.' {
				_, ok := digitSet[char]
				if !ok && !visits[i][j] {
					// neither number nor dot
					neighbors := getNeighborNumbers(i, j, lines, visits)
					collectedSymbols = append(collectedSymbols, struct {
						x         int
						y         int
						c         string
						neighbors []string
					}{x: i, y: j, c: string(char), neighbors: neighbors})

					if char == '*' && len(neighbors) > 1 {
						collectedGears = append(collectedGears, struct{ numbers []string }{neighbors})
					}
				}
			}
		}
	}

	sum := int64(0)
	for _, s := range collectedSymbols {
		for _, v := range s.neighbors {
			v, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				continue
			}
			sum += v
		}
	}

	fmt.Printf("Sum: %d\n", sum)

	gearRatiosSum := int64(0)
	for _, g := range collectedGears {
		ratio := int64(1)
		for _, n := range g.numbers {
			v, err := strconv.ParseInt(n, 10, 0)
			if err != nil {
				continue
			}
			ratio *= v
		}
		gearRatiosSum += ratio
	}
	fmt.Printf("Prod: %d\n", gearRatiosSum)
}
