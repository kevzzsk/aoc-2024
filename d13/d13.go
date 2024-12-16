package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"time"
)

type machineConfig struct {
	buttonA [2]int
	buttonB [2]int
	prize   [2]int
	cache   map[string]int
}

func extractNumbers(input string) [2]int {
	re := regexp.MustCompile(`[+=](\d+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	x, _ := strconv.Atoi(matches[0][1])
	y, _ := strconv.Atoi(matches[1][1])

	return [2]int{x, y}
}

// this function reads the input from the txt file and splits the rules and prints
func readInput(part2 bool) []machineConfig {
	// read from txt file
	file, err := os.Open("input13.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var machines []machineConfig

	// read the file line by line
	scanner := bufio.NewScanner(file)
	inputIndex := 0
	machine := machineConfig{
		cache: make(map[string]int),
	}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			inputIndex = 0
			continue
		}
		if inputIndex == 0 {
			// button A
			machine.buttonA = extractNumbers(line)
		}
		if inputIndex == 1 {
			// button B
			machine.buttonB = extractNumbers(line)
		}
		if inputIndex == 2 {
			// prize
			machine.prize = extractNumbers(line)
			if part2 {
				machine.prize[0] += 10000000000000
				machine.prize[1] += 10000000000000
			}
			machines = append(machines, machine)
			machine = machineConfig{
				cache: make(map[string]int),
			}
		}
		inputIndex++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return machines
}

func main() {
	startTime := time.Now()
	machines := readInput(false)

	minTokenToSpend := 0
	for _, machine := range machines {
		findMinToken := machine.findMinToken(machine.prize[0], machine.prize[1], 0, 0)
		if findMinToken != math.MaxInt {
			minTokenToSpend += findMinToken
		}

	}

	minTokenToSpend2 := 0
	machines2 := readInput(true)
	for _, machine := range machines2 {
		findMinToken := machine.findMinTokenWithLinearAlgebra()
		if findMinToken != 0 {
			minTokenToSpend2 += findMinToken
		}
	}

	fmt.Println("Total min token to spend:", minTokenToSpend)
	fmt.Println("Total min token to spend with linear algebra:", minTokenToSpend2)

	endTime := time.Since(startTime)
	fmt.Printf("Day13 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day13 execution took: 996 ms (996255 µs)
}

func (m *machineConfig) findMinTokenWithLinearAlgebra() int {
	determinant := m.buttonA[0]*m.buttonB[1] - m.buttonB[0]*m.buttonA[1]
	aPressed := (m.prize[0]*m.buttonB[1] - m.prize[1]*m.buttonB[0]) / determinant
	bPressed := (m.prize[1]*m.buttonA[0] - m.prize[0]*m.buttonA[1]) / determinant

	if (m.buttonA[0]*aPressed+m.buttonB[0]*bPressed) == m.prize[0] && (m.buttonA[1]*aPressed+m.buttonB[1]*bPressed) == m.prize[1] {
		return aPressed*3 + bPressed
	} else {
		return 0
	}

}

func (m *machineConfig) findMinToken(tx, ty int, aPressed, bPressed int) int {
	// base case
	if tx == 0 && ty == 0 {
		return aPressed*3 + bPressed
	}

	if tx < 0 || ty < 0 {
		return math.MaxInt
	}
	// cannot have more than 100 presses
	if aPressed > 100 || bPressed > 100 {
		return math.MaxInt
	}

	if value, exists := m.cache[fmt.Sprintf("%d,%d", tx, ty)]; exists {
		return value
	}

	minTokens := min(m.findMinToken(tx-m.buttonA[0], ty-m.buttonA[1], aPressed+1, bPressed), m.findMinToken(tx-m.buttonB[0], ty-m.buttonB[1], aPressed, bPressed+1))
	m.cache[fmt.Sprintf("%d,%d", tx, ty)] = minTokens
	return minTokens
}
