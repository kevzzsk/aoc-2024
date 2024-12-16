package main

import (
	utils "aoc-2024"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// this function reads the input from the txt file and splits the rules and prints
func readInput() [][]string {
	// read from txt file
	file, err := os.Open("input12.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var grid [][]string

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		grid = append(grid, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return grid
}

type garden struct {
	grid    [][]string
	visited map[string]bool
}

func (g *garden) getSameRegionNeighbors(x, y int) (up, right, down, left bool) {
	currentRegion := g.grid[y][x]
	if y > 0 && g.grid[y-1][x] == currentRegion {
		up = true
	}
	if x < len(g.grid[y])-1 && g.grid[y][x+1] == currentRegion {
		right = true
	}
	if y < len(g.grid)-1 && g.grid[y+1][x] == currentRegion {
		down = true
	}
	if x > 0 && g.grid[y][x-1] == currentRegion {
		left = true
	}
	return
}

func main() {
	startTime := time.Now()
	grid := readInput()

	garden := garden{
		grid:    grid,
		visited: make(map[string]bool),
	}

	price := 0
	discountedPrice := 0
	for iy, y := range grid {
		for ix, _ := range y {
			if !garden.visited[fmt.Sprintf("%d,%d", ix, iy)] {
				area, perimeter, corner := garden.traverseRegion(ix, iy)
				price += area * perimeter
				discountedPrice += area * corner
			}
		}
	}

	fmt.Println("Price:", price)
	fmt.Println("Discounted Price:", discountedPrice)
	endTime := time.Since(startTime)
	fmt.Printf("Day12 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day12 execution took: 16 ms (16519 µs)
}

// used in part2 to get the number of straight fences. As number of straight fences = number of corners
func (g *garden) getCurrCorner(x, y int) int {
	cornerCount := 0
	up, right, down, left := g.getSameRegionNeighbors(x, y)

	// check top left
	if !up && !left {
		cornerCount++
	}
	// check top right
	if !up && !right {
		cornerCount++
	}
	// check bottom right
	if !down && !right {
		cornerCount++
	}
	// check bottom left
	if !down && !left {
		cornerCount++
	}

	// check for inner corners
	// check top left
	if up && left && g.grid[y-1][x-1] != g.grid[y][x] {
		cornerCount++
	}
	// check top right
	if up && right && g.grid[y-1][x+1] != g.grid[y][x] {
		cornerCount++
	}
	// check bottom right
	if down && right && g.grid[y+1][x+1] != g.grid[y][x] {
		cornerCount++
	}
	// check bottom left
	if down && left && g.grid[y+1][x-1] != g.grid[y][x] {
		cornerCount++
	}

	return cornerCount
}

func (g *garden) traverseRegion(x, y int) (area, perimeter, corner int) {
	up, right, down, left := g.getSameRegionNeighbors(x, y)
	connectedNeighbors := utils.CountBool(up, right, down, left)

	g.visited[fmt.Sprintf("%d,%d", x, y)] = true

	currCorner := g.getCurrCorner(x, y)

	if up && !g.visited[fmt.Sprintf("%d,%d", x, y-1)] {
		neighboursArea, neighboursPerimeter, neighboursCorner := g.traverseRegion(x, y-1)
		area += neighboursArea
		perimeter += neighboursPerimeter
		corner += neighboursCorner
	}
	if right && !g.visited[fmt.Sprintf("%d,%d", x+1, y)] {
		neighboursArea, neighboursPerimeter, neighboursCorner := g.traverseRegion(x+1, y)
		area += neighboursArea
		perimeter += neighboursPerimeter
		corner += neighboursCorner
	}
	if down && !g.visited[fmt.Sprintf("%d,%d", x, y+1)] {
		neighboursArea, neighboursPerimeter, neighboursCorner := g.traverseRegion(x, y+1)
		area += neighboursArea
		perimeter += neighboursPerimeter
		corner += neighboursCorner
	}
	if left && !g.visited[fmt.Sprintf("%d,%d", x-1, y)] {
		neighboursArea, neighboursPerimeter, neighboursCorner := g.traverseRegion(x-1, y)
		area += neighboursArea
		perimeter += neighboursPerimeter
		corner += neighboursCorner
	}
	return area + 1, perimeter + (4 - connectedNeighbors), corner + currCorner
}
