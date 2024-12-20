package main

import (
	utils "aoc-2024"
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type point struct {
	x int
	y int
}

type maze struct {
	grid       [][]string
	startPoint point
	endPoint   point
	obstacles  []point
}

type Vector struct {
	point          point
	heuristicScore int // h(n) - estimated cost to reach the end node
	currentCost    int // g(n) - cost to reach the current node
	totalCost      int // f(n) - total cost of the node
	parent         *Vector
	step           int
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() maze {
	// read from txt file
	file, err := os.Open("input18.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// init grid with "." to represent empty space
	width := 71
	height := 71
	grid := make([][]string, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]string, width)
		for j := 0; j < width; j++ {
			grid[i][j] = "."
		}
	}

	maze := maze{
		grid:       grid,
		startPoint: point{0, 0},
	}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	maxFall := 1024
	for scanner.Scan() {
		line := scanner.Text()
		block := strings.Split(line, ",")
		blockX, _ := strconv.Atoi(block[0])
		blockY, _ := strconv.Atoi(block[1])
		if maxFall > 0 {
			maze.grid[blockY][blockX] = "#"
			maxFall--
		} else {
			maze.obstacles = append(maze.obstacles, point{blockX, blockY})
		}
	}
	maze.endPoint = point{len(maze.grid[0]) - 1, len(maze.grid) - 1}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return maze
}

func (m *maze) print() {

	for iy, row := range m.grid {
		for ix, v := range row {
			if ix == m.startPoint.x && iy == m.startPoint.y {
				fmt.Print("S")
			} else if ix == m.endPoint.x && iy == m.endPoint.y {
				fmt.Print("E")
			} else {
				fmt.Print(v)
			}
		}
		fmt.Println()
	}
}

func main() {
	startTime := time.Now()

	m := readInput()

	fmt.Println("steps taken: ", m.traverse())

	for _, obstacle := range m.obstacles {
		// add in obstacle and try to find if there is a path
		m.grid[obstacle.y][obstacle.x] = "#"
		res := m.traverse()
		if res == -1 {
			// no path is found
			fmt.Println("obstacle causing no path: ", obstacle)
			break
		}
	}

	endTime := time.Since(startTime)
	fmt.Printf("Day18 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day18 execution took: 24991 ms (24991117 µs)
}

// use the manhattan distance which measures the distance between two points in a grid based an straight line
func (m *maze) heuristicScore(p point) int {
	return utils.Abs(p.x-m.endPoint.x) + utils.Abs(p.y-m.endPoint.y)
}

func (m *maze) costToMove(current, target point) int {
	// check if we are at the same point
	if current.x == target.x && current.y == target.y {
		return 0
	}

	return 1
}

func (m *maze) getAllNeighbouringPathPoints(p point) []point {
	neighbours := []point{
		{p.x, p.y - 1}, // N
		{p.x + 1, p.y}, // E
		{p.x, p.y + 1}, // S
		{p.x - 1, p.y}, // W
	}
	validNeighbours := []point{}
	// check if the point is a wall or out of bounds
	for i, v := range neighbours {
		if v.x < 0 || v.x >= len(m.grid[0]) || v.y < 0 || v.y >= len(m.grid) {
			continue
		}
		if m.grid[v.y][v.x] != "#" {
			validNeighbours = append(validNeighbours, neighbours[i])
		}
	}

	return validNeighbours
}

// todo: update to be more efficient
func ItemExists(pq *utils.PriorityQueue[Vector], p point) (*utils.Item[Vector], error) {
	for _, v := range *pq {
		if reflect.DeepEqual(v.Value.point, p) {
			return v, nil
		}
	}
	return nil, fmt.Errorf("item does not exist in the priority queue")
}

func (m *maze) traverse() int {
	pq := make(utils.PriorityQueue[Vector], 0)
	heap.Init(&pq)

	visited := make(map[point]bool)

	startHeuristicScore := m.heuristicScore(m.startPoint)
	startTotalCost := startHeuristicScore
	startItem := &utils.Item[Vector]{
		Value: Vector{
			point:          m.startPoint,
			heuristicScore: startHeuristicScore,
			currentCost:    0,
			totalCost:      startTotalCost,
			parent:         nil,
			step:           0,
		},
		Priority: startTotalCost,
	}
	heap.Push(&pq, startItem)

	for pq.Len() > 0 {
		currentVector := heap.Pop(&pq).(*utils.Item[Vector]).Value

		// check if we reached the end point
		if currentVector.point.x == m.endPoint.x && currentVector.point.y == m.endPoint.y {
			return currentVector.totalCost

		}

		// add current point to visited
		visited[currentVector.point] = true

		// check all possible directions to move
		for _, point := range m.getAllNeighbouringPathPoints(currentVector.point) {

			// skip if point already visited
			if visited[point] {
				continue
			}

			estimatedCostToBeInPoint := m.costToMove(currentVector.point, point) + currentVector.currentCost

			heuristicScore := m.heuristicScore(point)
			currentCost := estimatedCostToBeInPoint
			totalCost := estimatedCostToBeInPoint + heuristicScore
			neighbourItem := &utils.Item[Vector]{
				Value: Vector{
					point:          point,
					heuristicScore: heuristicScore,
					currentCost:    currentCost,
					totalCost:      totalCost,
					parent:         &currentVector,
					step:           currentVector.step + 1,
				},
				Priority: totalCost,
			}

			existingNeighbourItem, _ := ItemExists(&pq, point)
			if existingNeighbourItem == nil {
				// add to queue
				heap.Push(&pq, neighbourItem)
			} else {

				// check if neighbour is a better path
				if neighbourItem.Value.totalCost <= totalCost {
					continue
				} else {
					// we are currently on better path - update neighbour
					pq.Update(neighbourItem, neighbourItem.Value, neighbourItem.Priority)
				}
			}
		}
	}
	return -1
}
