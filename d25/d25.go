package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type schematics struct {
	keys [][][]string
	lock [][][]string

	pkeys [][]int
	plock [][]int
}

func readInput() schematics {
	// read from txt file
	file, err := os.Open("input25.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := schematics{}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	index := 0
	islock := false
	currentSchematic := [][]string{}
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if line == "" {
			index = 0
			currentSchematic = [][]string{}
			islock = false
			continue
		}
		if index == 0 && line == "#####" {
			islock = true
		}
		currentSchematic = append(currentSchematic, strings.Split(line, ""))
		if index == 6 {
			if !islock {
				result.keys = append(result.keys, currentSchematic)
			} else {
				result.lock = append(result.lock, currentSchematic)
			}
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func main() {
	startTime := time.Now()

	schema := readInput()

	for _, key := range schema.keys {
		numbers := make([]int, 5)
		for _, row := range key[:6] {
			for j, col := range row {
				if col == "#" {
					numbers[j] += 1
				}
			}
		}
		schema.pkeys = append(schema.pkeys, numbers)
	}

	for _, lock := range schema.lock {
		numbers := make([]int, 5)
		for _, row := range lock[1:] {
			for j, col := range row {
				if col == "#" {
					numbers[j] += 1
				}
			}
		}
		schema.plock = append(schema.plock, numbers)
	}

	uniqueLock := 0
	for _, key := range schema.pkeys {
		// check if same as any lock
		for _, lock := range schema.plock {
			if isLockAndKeyFit(key, lock) {
				uniqueLock++
			}
		}
	}

	fmt.Println("uniqueLock:", uniqueLock)

	endTime := time.Since(startTime)
	fmt.Printf("Day25 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day25 execution took: 32 ms (32867 µs)
}

func isLockAndKeyFit(key, lock []int) bool {
	for i := 0; i < 5; i++ {
		if key[i]+lock[i] > 5 {
			return false
		}
	}
	return true
}
