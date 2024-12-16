package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// this function reads the input from the txt file and splits the rules and prints
func readInput() []int {
	// read from txt file
	file, err := os.Open("input11.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var stones []int

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		for _, v := range strings.Split(line, " ") {
			stone, _ := strconv.Atoi(v)
			stones = append(stones, stone)
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return stones
}

type blinkStone struct {
	startingStone int
	blinkCount    int
	stoneCount    int
	stones        []int
}

func main() {
	startTime := time.Now()
	stones := readInput()

	fmt.Println(stones)

	blinkDict := make(map[string]int)
	fmt.Println(blinkBFS(125, 25, &blinkDict) + blinkBFS(17, 25, &blinkDict))

	stonesCount25 := 0
	for _, stone := range stones {
		stonesCount25 += blinkBFS(stone, 25, &blinkDict)
	}

	stonesCount75 := 0
	for _, stone := range stones {
		stonesCount75 += blinkBFS(stone, 75, &blinkDict)
	}

	fmt.Println()
	fmt.Println("StonesCount25:", stonesCount25)
	fmt.Println("StonesCount75:", stonesCount75)
	fmt.Println(len(blinkDict))
	endTime := time.Since(startTime)
	fmt.Printf("Day11 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day11
}

func printStoneBlinkDict(stoneBlinkDict map[int]map[int]blinkStone) {
	for k, v := range stoneBlinkDict {
		fmt.Println("Stone:", k)
		for k1, v1 := range v {
			fmt.Printf("stone: %d, blinkCount: %d, Value: %v\n", k, k1, v1)
		}
	}
}

func (sbd *blinkStone) String() string {
	return fmt.Sprintf("StartingStone: %d, BlinkCount: %d, StoneCount: %d, Stones: %v", sbd.startingStone, sbd.blinkCount, sbd.stoneCount, sbd.stones)
}

// split even length number down the middle
// e.g. 1221 -> 12, 21
func splitEvenNumbers(n int) []int {
	str := strconv.Itoa(n)
	mid := len(str) / 2
	firstHalf, _ := strconv.Atoi(str[:mid])
	secondHalf, _ := strconv.Atoi(str[mid:])
	return []int{firstHalf, secondHalf}
}

func hasEvenLength(stone int) bool {
	return len(strconv.Itoa(stone))%2 == 0
}

func blink(stone int, sbd *map[int]map[int]blinkStone) []int {
	if _, exists := (*sbd)[stone]; exists {
		if blinkStone, exists := (*sbd)[stone][1]; exists {
			return blinkStone.stones
		}
	}

	if stone == 0 {
		return []int{1}
	} else if hasEvenLength(stone) {
		return splitEvenNumbers(stone)
	}
	return []int{stone * 2024}
}

func blink2(stone int) []int {
	if stone == 0 {
		return []int{1}
	} else if hasEvenLength(stone) {
		return splitEvenNumbers(stone)
	}
	return []int{stone * 2024}
}

func blinkBFS(stone int, blinkCount int, blinkDict *map[string]int) int {
	if blinkCount == 0 {
		return 1
	}
	if val, exists := (*blinkDict)[fmt.Sprintf("%d,%d", stone, blinkCount)]; exists {
		return val
	}
	count := 0

	if stone == 0 {
		count = blinkBFS(1, blinkCount-1, blinkDict)
	} else if hasEvenLength(stone) {
		splitStones := splitEvenNumbers(stone)
		count = blinkBFS(splitStones[0], blinkCount-1, blinkDict) + blinkBFS(splitStones[1], blinkCount-1, blinkDict)
	} else {
		count = blinkBFS(stone*2024, blinkCount-1, blinkDict)
	}
	(*blinkDict)[fmt.Sprintf("%d,%d", stone, blinkCount)] = count
	return count

}
