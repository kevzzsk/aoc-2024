package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type topographicMap struct {
	grid [][]int
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() [][]int {
	// read from txt file
	file, err := os.Open("input10.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]int

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, v := range strings.Split(line, "") {
			row[i], _ = strconv.Atoi(v)
		}
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return grid
}

func main() {
	startTime := time.Now()
	grid := readInput()

	tm := topographicMap{grid: grid}

	const TRAILHEAD = 0

	sumOfScores := 0
	sumOfRating := 0
	for iy, row := range tm.grid {
		for ix, v := range row {
			if v == TRAILHEAD {
				// PART 1
				visited := map[string]bool{}
				countHighestPointReachable(tm, ix, iy, 0, visited)
				trailheadScore := len(visited)
				sumOfScores += trailheadScore

				// PART 2
				rating := countPathToAllPeak(tm, ix, iy, 0)
				sumOfRating += rating
			}
		}
	}

	fmt.Println("Sum of scores: ", sumOfScores)
	fmt.Println("Sum of rating: ", sumOfRating)

	endTime := time.Since(startTime)
	fmt.Printf("Day10 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day10 execution took: 1 ms (1280 µs)
}

// return -1 if the point is not in map
func (tm topographicMap) getNeighbours(x, y int) (up, right, down, left int) {
	up, right, down, left = -1, -1, -1, -1
	if y > 0 {
		up = tm.grid[y-1][x]
	}
	if x < len(tm.grid[y])-1 {
		right = tm.grid[y][x+1]
	}
	if y < len(tm.grid)-1 {
		down = tm.grid[y+1][x]
	}
	if x > 0 {
		left = tm.grid[y][x-1]
	}

	return up, right, down, left
}

func countHighestPointReachable(tm topographicMap, x, y, level int, visited map[string]bool) {
	if level == 9 {
		visited[fmt.Sprintf("%d,%d", x, y)] = true
		return
	}
	up, right, down, left := tm.getNeighbours(x, y)

	if up == level+1 {
		countHighestPointReachable(tm, x, y-1, level+1, visited)
	}
	if right == level+1 {
		countHighestPointReachable(tm, x+1, y, level+1, visited)
	}
	if down == level+1 {
		countHighestPointReachable(tm, x, y+1, level+1, visited)
	}
	if left == level+1 {
		countHighestPointReachable(tm, x-1, y, level+1, visited)
	}
}

func countPathToAllPeak(tm topographicMap, x, y, level int) int {
	if level == 9 {
		return 1
	}
	up, right, down, left := tm.getNeighbours(x, y)
	count := 0
	if up == level+1 {
		count += countPathToAllPeak(tm, x, y-1, level+1)
	}
	if right == level+1 {
		count += countPathToAllPeak(tm, x+1, y, level+1)
	}
	if down == level+1 {
		count += countPathToAllPeak(tm, x, y+1, level+1)
	}
	if left == level+1 {
		count += countPathToAllPeak(tm, x-1, y, level+1)
	}
	return count
}
