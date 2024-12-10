package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Board struct {
	grid             [][]string
	currentXPos      int
	currentYPos      int
	currentDirection Direction
	visited          map[string]bool // keys stored as location visited like "x,y"
	visitedVector    map[string]bool // keys stored as location visited like "x,y,direction" -> used for optimization in part 2

}

type Direction int

const (
	N Direction = iota
	E
	S
	W
)

func (b *Board) findCurrentPos() {
	const GUARD = "^"

	for iy, y := range b.grid {
		for ix, x := range y {
			if x == GUARD {
				b.currentXPos = ix
				b.currentYPos = iy
			}
		}
	}
}

func (b *Board) isXWithinGrid(x int) bool {
	if x < 0 || x >= len(b.grid[0]) {
		return false
	}
	return true
}

func (b *Board) isYWithinGrid(y int) bool {
	if y < 0 || y >= len(b.grid) {
		return false
	}
	return true
}

func (b *Board) getForwardPosition(x, y int, direction Direction) (int, int) {
	switch direction {
	case N:
		y = y - 1
	case E:
		x = x + 1
	case S:
		y = y + 1
	case W:
		x = x - 1
	}
	return x, y
}

// peekForward returns the value of the next position in the direction
// if the next position is out of bounds, it returns ""
// if the next position is an obstacle, it returns "#"
func (b *Board) peekForward(x, y int, direction Direction) string {
	x, y = b.getForwardPosition(x, y, direction)

	// check if the next position is out of bounds => return ""
	if !b.isXWithinGrid(x) || !b.isYWithinGrid(y) {
		return ""
	}

	return b.grid[y][x]
}

func (b *Board) printCurrentGrid() {
	// print the grid
	for iy, y := range b.grid {
		for ix, x := range y {
			if b.currentXPos == ix && b.currentYPos == iy {
				// print ^ if current direction is N
				if b.currentDirection == N {
					fmt.Print("^")
				} else if b.currentDirection == E {
					fmt.Print(">")
				} else if b.currentDirection == S {
					fmt.Print("v")
				} else {
					fmt.Print("<")
				}
			} else if b.visited[fmt.Sprintf("%d,%d", ix, iy)] {
				// if visited print X
				fmt.Print("X")
			} else {
				fmt.Print(x)
			}
		}
		fmt.Println()
	}
}

// check if putting a blockade in front of the current path causes loop
// simple approach is to check if the right side of the current path is visited
// if visited, we can put a blockade infront to force guard to go right and go into a loop
func (b *Board) checkIfBlockCausesLoop() bool {

	ghostXPos := b.currentXPos
	ghostYPos := b.currentYPos
	ghostDirection := b.currentDirection
	ghostVisited := make(map[string]bool)

	// walk until edge of map or loop is detected
	for ghostItemInfront := b.peekForward(ghostXPos, ghostYPos, ghostDirection); ; ghostItemInfront = b.peekForward(ghostXPos, ghostYPos, ghostDirection) {
		if ghostItemInfront == "" {
			return false
		}
		if ghostItemInfront == "#" {
			// turn right (90 degrees)
			ghostDirection = (ghostDirection + 1) % 4

		} else {
			ghostXPos, ghostYPos = b.getForwardPosition(ghostXPos, ghostYPos, ghostDirection)
		}
		// detect loop
		if ghostVisited[fmt.Sprintf("%d,%d,%d", ghostXPos, ghostYPos, ghostDirection)] || b.visitedVector[fmt.Sprintf("%d,%d,%d", ghostXPos, ghostYPos, ghostDirection)] {
			return true
		}
		ghostVisited[fmt.Sprintf("%d,%d,%d", ghostXPos, ghostYPos, ghostDirection)] = true

	}
}

func (b *Board) traverse() {
	b.currentDirection = N
	// include the starting position as visited
	b.visited[fmt.Sprintf("%d,%d", b.currentXPos, b.currentYPos)] = true
	b.visitedVector[fmt.Sprintf("%d,%d,%d", b.currentXPos, b.currentYPos, b.currentDirection)] = true

	ghostObstacleMap := make(map[string]bool)

	for objectInfront := b.peekForward(b.currentXPos, b.currentYPos, b.currentDirection); objectInfront != ""; objectInfront = b.peekForward(b.currentXPos, b.currentYPos, b.currentDirection) {
		// for part 2 - check if putting a blockade in front of the current path causes loop
		ghostObstacleXPos, ghostObstacleYPos := b.getForwardPosition(b.currentXPos, b.currentYPos, b.currentDirection)
		// put temporary blockade and check if it causes loop
		temp := b.grid[ghostObstacleYPos][ghostObstacleXPos]
		b.grid[ghostObstacleYPos][ghostObstacleXPos] = "#"
		if !b.visited[fmt.Sprintf("%d,%d", ghostObstacleXPos, ghostObstacleYPos)] && b.checkIfBlockCausesLoop() {
			ghostObstacleMap[fmt.Sprintf("%d,%d", ghostObstacleXPos, ghostObstacleYPos)] = true

		}
		b.grid[ghostObstacleYPos][ghostObstacleXPos] = temp

		// walk towards the current direction until we hit a wall
		if objectInfront == "#" {
			// turn right (90 degrees)
			b.currentDirection = (b.currentDirection + 1) % 4
		} else {
			// walk 1 step and record the location
			b.currentXPos, b.currentYPos = b.getForwardPosition(b.currentXPos, b.currentYPos, b.currentDirection)
		}
		b.visited[fmt.Sprintf("%d,%d", b.currentXPos, b.currentYPos)] = true
		b.visitedVector[fmt.Sprintf("%d,%d,%d", b.currentXPos, b.currentYPos, b.currentDirection)] = true
	}

	// guard exits the map => count distinct locations visited
	fmt.Println("Total distinct locations visited: ", len(b.visited))
	fmt.Println("Loop blockade count: ", len(ghostObstacleMap))
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() [][]string {
	// read from txt file
	file, err := os.Open("input6.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var board [][]string

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		board = append(board, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return board
}

func main() {
	startTime := time.Now()
	board := readInput()

	b := Board{
		grid:          board,
		visited:       make(map[string]bool),
		visitedVector: make(map[string]bool),
	}

	b.findCurrentPos()
	b.traverse()

	endTime := time.Since(startTime)
	fmt.Printf("Day6 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) //Day6 execution took: 2973 ms (2973427 µs)
}
