package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Game struct {
	red   int64
	green int64
	blue  int64
}

type GameSets map[int64][]Game

type PossibleGames map[int64]bool

func Max(array []int64) int64 {
	max_ := array[0]
	for _, value := range array {
		if max_ < value {
			max_ = value
		}
	}
	return max_
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

	gameRegExp, err := regexp.Compile(`^Game (\d+)`)
	redRegExp, err := regexp.Compile(`(\d+) red`)
	greenRegExp, err := regexp.Compile(`(\d+) green`)
	blueRegExp, err := regexp.Compile(`(\d+) blue`)
	if err != nil {
		fmt.Println("could not compile regular expression")
		os.Exit(1)
		return
	}

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	gameSets := make(GameSets)

	for scanner.Scan() {
		line := scanner.Text()
		split := strings.Split(line, ": ")

		if len(split) != 2 {
			continue
		}

		gameNumber := gameRegExp.FindStringSubmatch(split[0])
		if gameNumber == nil || len(gameNumber) < 2 {
			continue
		}

		game, err := strconv.ParseInt(string(gameNumber[1]), 10, 0)
		if err != nil {
			continue
		}
		sets := strings.Split(split[1], "; ")

		g := make([]Game, 0)
		for _, set := range sets {
			s := Game{red: 0, green: 0, blue: 0}
			redNumber := redRegExp.FindStringSubmatch(set)
			if redNumber != nil && len(redNumber) == 2 {
				red, rerr := strconv.ParseInt(string(redNumber[1]), 10, 0)
				if rerr == nil {
					s.red = red
				}
			}

			greenNumber := greenRegExp.FindStringSubmatch(set)
			if greenNumber != nil && len(greenNumber) == 2 {
				green, err := strconv.ParseInt(string(greenNumber[1]), 10, 0)
				if err == nil {
					s.green = green
				}
			}

			blueNumber := blueRegExp.FindStringSubmatch(set)
			if blueNumber != nil && len(blueNumber) == 2 {
				blue, err := strconv.ParseInt(string(blueNumber[1]), 10, 0)
				if err == nil {
					s.blue = blue
				}
			}
			g = append(g, s)
		}
		gameSets[game] = g
	}

	possibleGames := make(PossibleGames, 0)
	fewestPerGame := make(map[int64]Game)

	for gameID, sets := range gameSets {

		reds := make([]int64, len(sets))
		greens := make([]int64, len(sets))
		blues := make([]int64, len(sets))

		for _, game := range sets {
			reds = append(reds, game.red)
			greens = append(greens, game.green)
			blues = append(blues, game.blue)
		}

		fewestPerGame[gameID] = Game{
			red:   Max(reds),
			green: Max(greens),
			blue:  Max(blues),
		}

		for _, game := range sets {
			if possible, ok := possibleGames[gameID]; ok && !possible {
				break // game exists and is already marked as impossible
			}
			if game.red > 12 || game.green > 13 || game.blue > 14 {
				possibleGames[gameID] = false
				continue
			} else {
				possibleGames[gameID] = true
			}
		}
	}

	sum := int64(0)
	for gameID, possible := range possibleGames {
		if possible {
			sum += gameID
		}
	}

	fmt.Println("Sum: ", sum)

	power := int64(0)
	for _, g := range fewestPerGame {
		power += g.red * g.green * g.blue
	}
	fmt.Println("power: ", power)
}
