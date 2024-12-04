package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func readInput() [][]string {
	// read from txt file
	file, err := os.Open("input4.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data [][]string

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		letters := strings.Split(line, "")
		row := []string{}
		row = append(row, letters...)
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return data
}

func countTrue(b ...bool) int {
	n := 0
	for _, isTrue := range b {
		if isTrue {
			n++
		}
	}
	return n
}

func main() {
	startTime := time.Now()
	data := readInput()

	p1(&data)
	p2(&data)

	endTime := time.Since(startTime)
	fmt.Printf("Day4 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds())
}

// Find XMAS horizontally, vertically, and diagonally
// approach is to find X and check surrounding for M, A, S
// Time complexity: O(n * m) where n is the number of rows and m is the number of columns
// Space complexity: O(1)
func p1(data *[][]string) {
	xmasCount := 0

	for i, row := range *data {
		for j, item := range row {
			if item == "X" {
				// check surrounding for XMAS
				xmasCount += countTrue(
					checkWest(i, j, data),
					checkEast(i, j, data),
					checkNorth(i, j, data),
					checkSouth(i, j, data),
					checkNorthWest(i, j, data),
					checkNorthEast(i, j, data),
					checkSouthWest(i, j, data),
					checkSouthEast(i, j, data))
			}
		}
	}

	fmt.Println("XMAS Count: ", xmasCount)
}

// SAMX
func checkWest(i int, j int, data *[][]string) bool {
	if j < 3 {
		return false
	}
	return (*data)[i][j-1] == "M" && (*data)[i][j-2] == "A" && (*data)[i][j-3] == "S"
}

// XMAS
func checkEast(i int, j int, data *[][]string) bool {
	if j > len((*data)[0])-1-3 {
		return false
	}
	return (*data)[i][j+1] == "M" && (*data)[i][j+2] == "A" && (*data)[i][j+3] == "S"
}

// S
// A
// M
// X
func checkNorth(i int, j int, data *[][]string) bool {
	if i < 3 {
		return false
	}
	return (*data)[i-1][j] == "M" && (*data)[i-2][j] == "A" && (*data)[i-3][j] == "S"
}

// X
// M
// A
// S
func checkSouth(i int, j int, data *[][]string) bool {
	if i > len((*data))-1-3 {
		return false
	}
	return (*data)[i+1][j] == "M" && (*data)[i+2][j] == "A" && (*data)[i+3][j] == "S"
}

// S...
// .M..
// ..A.
// ...X
func checkNorthWest(i int, j int, data *[][]string) bool {
	if i < 3 || j < 3 {
		return false
	}
	return (*data)[i-1][j-1] == "M" && (*data)[i-2][j-2] == "A" && (*data)[i-3][j-3] == "S"
}

// ...S
// ..A.
// .M..
// X...
func checkNorthEast(i int, j int, data *[][]string) bool {
	if i < 3 || j > len((*data)[0])-1-3 {
		return false
	}
	return (*data)[i-1][j+1] == "M" && (*data)[i-2][j+2] == "A" && (*data)[i-3][j+3] == "S"
}

// ...X
// ..M.
// .A..
// S...
func checkSouthWest(i int, j int, data *[][]string) bool {
	if i > len((*data))-1-3 || j < 3 {
		return false
	}
	return (*data)[i+1][j-1] == "M" && (*data)[i+2][j-2] == "A" && (*data)[i+3][j-3] == "S"
}

// X...
// .M..
// ..A.
// ...S
func checkSouthEast(i int, j int, data *[][]string) bool {
	if i > len((*data))-1-3 || j > len((*data)[0])-1-3 {
		return false
	}
	return (*data)[i+1][j+1] == "M" && (*data)[i+2][j+2] == "A" && (*data)[i+3][j+3] == "S"
}

// find 2 MAS crossing each other diagonally
// approach is to find A and check surrounding for M, A, S
// Time complexity: O(n * m) where n is the number of rows and m is the number of columns
// Space complexity: O(1)
func p2(data *[][]string) {
	xmasCount := 0

	for i, row := range *data {
		for j, item := range row {
			if item == "A" {
				// check surrounding for X-MAS
				if exists := checkXMAS(i, j, data); exists {
					xmasCount++
				}
			}
		}
	}

	fmt.Println("X-MAS Count: ", xmasCount)
}

func checkXMAS(i int, j int, data *[][]string) bool {
	if i == 0 || i == len((*data))-1 || j == 0 || j == len((*data)[0])-1 {
		// if A is in corner, not possible for X-MAS
		return false
	}

	topLeft := (*data)[i+1][j-1]
	topRight := (*data)[i+1][j+1]
	botLeft := (*data)[i-1][j-1]
	botRight := (*data)[i-1][j+1]

	// CASE 1: M.S
	//			A
	//		   M.S
	if topLeft == "M" && topRight == "S" && botLeft == "M" && botRight == "S" {
		return true
	}
	// CASE 2: M.M
	//			A
	//		   S.S
	if topLeft == "M" && topRight == "M" && botLeft == "S" && botRight == "S" {
		return true
	}
	// CASE 3: S.M
	//			A
	//		   S.M
	if topLeft == "S" && topRight == "M" && botLeft == "S" && botRight == "M" {
		return true
	}
	// CASE 4: S.S
	//			A
	//		   M.M
	if topLeft == "S" && topRight == "S" && botLeft == "M" && botRight == "M" {
		return true
	}

	return false
}
