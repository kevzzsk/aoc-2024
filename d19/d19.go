package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type towels struct {
	patterns []string
	designs  []string
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() towels {
	// read from txt file
	file, err := os.Open("input19.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	towels := towels{}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	readingPatterns := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingPatterns = false
			continue
		}
		if readingPatterns {
			row := strings.Split(line, ",")
			// trim the spaces
			for _, v := range row {
				towels.patterns = append(towels.patterns, strings.TrimSpace(v))
			}
		} else {
			towels.designs = append(towels.designs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return towels
}

func main() {
	startTime := time.Now()

	towel := readInput()

	possibleDesign := 0
	cache := make(map[string]bool)
	for _, design := range towel.designs {
		// check if any of the pattern can build the design
		if dfs(towel.patterns, design, cache) {
			possibleDesign++
		}
	}
	fmt.Println("Possible designs: ", possibleDesign)

	totalPermutationOfDesign := 0
	cache2 := make(map[string]int)
	for _, design := range towel.designs {
		// check if any of the pattern can build the design
		totalPermutationOfDesign += dfs2(towel.patterns, design, cache2)
	}

	fmt.Println("Total permutation of designs: ", totalPermutationOfDesign)

	endTime := time.Since(startTime)
	fmt.Printf("Day19 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day19 execution took: 56 ms (56256 µs)
}

func dfs(patterns []string, design string, cache map[string]bool) bool {
	if design == "" {
		return true
	}
	if isPossible, found := cache[design]; found {
		return isPossible
	}

	isPossible := false
	for _, pattern := range patterns {
		if newDesign, found := strings.CutPrefix(design, pattern); found {
			isPossible = isPossible || dfs(patterns, newDesign, cache)
		}
	}
	cache[design] = isPossible
	return isPossible
}

func dfs2(patterns []string, design string, cache map[string]int) int {
	if design == "" {
		return 1
	}
	if count, found := cache[design]; found {
		return count
	}

	count := 0
	for _, pattern := range patterns {
		if newDesign, found := strings.CutPrefix(design, pattern); found {
			count += dfs2(patterns, newDesign, cache)
		}
	}
	cache[design] = count
	return count
}
