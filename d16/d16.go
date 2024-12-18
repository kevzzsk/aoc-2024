package main

import (
	utils "aoc-2024"
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"time"
)

type point struct {
	x int
	y int
}

type pointv struct {
	p point
	d direction
}

type direction int

const (
	N direction = iota
	E
	S
	W
)

type maze struct {
	grid         [][]string
	startPoint   point
	endPoint     point
	bestPath     map[int]Vector
	bestPathCost int
	allBestPaths map[pointv]Vector
}

type Vector struct {
	point          point
	d              direction
	heuristicScore int // h(n) - estimated cost to reach the end node
	currentCost    int // g(n) - cost to reach the current node
	totalCost      int // f(n) - total cost of the node
	parent         *Vector
	step           int
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() maze {
	// read from txt file
	file, err := os.Open("input16.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	maze := maze{
		bestPath:     make(map[int]Vector),
		allBestPaths: make(map[pointv]Vector),
		bestPathCost: math.MaxInt,
	}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, "")
		for i, v := range row {
			if v == "S" {
				maze.startPoint = point{i, len(maze.grid)}
			} else if v == "E" {
				maze.endPoint = point{i, len(maze.grid)}
			}
		}
		maze.grid = append(maze.grid, row)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return maze
}

func (m *maze) print() {
	for _, row := range m.grid {
		fmt.Println(row)
	}
}

func (m *maze) printCurrentState(currV Vector, visited map[point]bool, pq utils.PriorityQueue[Vector]) {
	// convert all bestpath to remove the direction
	allBestPaths := make(map[point]bool)
	for k, _ := range m.allBestPaths {
		allBestPaths[k.p] = true
	}

	for iy, row := range m.grid {
		for ix, v := range row {
			point := point{ix, iy}
			if vector, _ := ItemExists(&pq, point); vector != nil {
				fmt.Print("P")
			} else if currV.point.x == ix && currV.point.y == iy {
				fmt.Print("C")
			} else if allBestPaths[point] {
				fmt.Print("A")
			} else if visited[point] {
				fmt.Print("X")
			} else {
				fmt.Print(v)
			}

		}
		fmt.Println()
	}
}

// Part 1 used A* algorithm to find the first best path
// Part 2 used modified A* to find all best paths -> not optimal. Should have used Dijkstra's algorithm
func main() {
	startTime := time.Now()
	m := readInput()

	m.traverse()
	fmt.Println("Best path cost: ", m.bestPathCost)
	for i := 0; i < 2; i++ {
		m.traverse2()
	}

	allBestPaths := make(map[point]bool)
	for k, _ := range m.allBestPaths {
		allBestPaths[k.p] = true
	}
	fmt.Println("All best path cost: ", len(allBestPaths)+1)

	endTime := time.Since(startTime)
	fmt.Printf("Day16 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day16 execution took: 1948 ms (1948631 µs)
}

// use the manhattan distance which measures the distance between two points in a grid based an straight line
func (m *maze) heuristicScore(p point) int {
	return utils.Abs(p.x-m.endPoint.x) + utils.Abs(p.y-m.endPoint.y)
}

// if we are moving forward in the same direction, cost is 1
// if we are moving in a different direction that requires turning, cost is 1000 per 90 degree turn
func (m *maze) costToMove(current, target point, currentDirection direction) int {
	// check if we are at the same point
	if current.x == target.x && current.y == target.y {
		return 0
	}

	// check if we are moving in the same direction
	if currentDirection == N {
		if current.y > target.y {
			return 1
		}
		if current.x != target.x {
			return 1001
		}
		if current.y < target.y {
			return 2001
		}
	}
	if currentDirection == E {
		if current.x < target.x {
			return 1
		}
		if current.y != target.y {
			return 1001
		}
		if current.x > target.x {
			return 2001
		}
	}
	if currentDirection == S {
		if current.y < target.y {
			return 1
		}
		if current.x != target.x {
			return 1001
		}
		if current.y > target.y {
			return 2001
		}
	}
	if currentDirection == W {
		if current.x > target.x {
			return 1
		}
		if current.y != target.y {
			return 1001
		}
		if current.x < target.x {
			return 2001
		}
	}
	return 0
}

func (m *maze) getNextDirection(current, target point, currentDirection direction) direction {
	if current.x == target.x && current.y == target.y {
		return currentDirection
	}

	if current.x == target.x {
		if current.y > target.y {
			return N
		}
		if current.y < target.y {
			return S
		}
	}
	if current.y == target.y {
		if current.x > target.x {
			return W
		}
		if current.x < target.x {
			return E
		}
	}
	return currentDirection

}

func (m *maze) getAllNeighbouringPathPoints(p point) []point {
	neighbours := []point{
		{p.x, p.y - 1}, // N
		{p.x + 1, p.y}, // E
		{p.x, p.y + 1}, // S
		{p.x - 1, p.y}, // W
	}
	validNeighbours := []point{}
	// check if the point is a wall
	for i, v := range neighbours {
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
			d:              E, // DEFAULT alwawys start facing east
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
			m.bestPathCost = currentVector.totalCost
			m.storeBestPath(currentVector)
			m.storeAllBestPath(currentVector)
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

			estimatedCostToBeInPoint := m.costToMove(currentVector.point, point, currentVector.d) + currentVector.currentCost

			heuristicScore := m.heuristicScore(point)
			currentCost := estimatedCostToBeInPoint
			totalCost := estimatedCostToBeInPoint + heuristicScore
			neighbourItem := &utils.Item[Vector]{
				Value: Vector{
					point:          point,
					d:              m.getNextDirection(currentVector.point, point, currentVector.d),
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

func (m *maze) storeBestPath(currentVector Vector) {
	for currentVector.parent != nil {
		m.bestPath[currentVector.step] = currentVector
		currentVector = *currentVector.parent
	}
}

func (m *maze) storeAllBestPath(currentVector Vector) {
	for currentVector.parent != nil {
		m.allBestPaths[pointv{currentVector.point, currentVector.d}] = currentVector
		currentVector = *currentVector.parent
	}
}

func (m *maze) traverse2() int {
	pq := make(utils.PriorityQueue[Vector], 0)
	heap.Init(&pq)

	visited := make(map[pointv]bool)

	startHeuristicScore := m.heuristicScore(m.startPoint)
	startTotalCost := startHeuristicScore
	startItem := &utils.Item[Vector]{
		Value: Vector{
			point:          m.startPoint,
			d:              E, // DEFAULT alwawys start facing east
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
			m.bestPathCost = currentVector.totalCost
			continue

		}

		// add current point to visited
		visited[pointv{currentVector.point, currentVector.d}] = true

		// check all possible directions to move
		for _, point := range m.getAllNeighbouringPathPoints(currentVector.point) {

			estimatedCostToBeInPoint := m.costToMove(currentVector.point, point, currentVector.d) + currentVector.currentCost

			heuristicScore := m.heuristicScore(point)
			currentCost := estimatedCostToBeInPoint
			totalCost := estimatedCostToBeInPoint + heuristicScore
			neighbourItem := &utils.Item[Vector]{
				Value: Vector{
					point:          point,
					d:              m.getNextDirection(currentVector.point, point, currentVector.d),
					heuristicScore: heuristicScore,
					currentCost:    currentCost,
					totalCost:      totalCost,
					parent:         &currentVector,
					step:           currentVector.step + 1,
				},
				Priority: totalCost,
			}

			//skip if point already visited
			if _, exists := m.allBestPaths[pointv{point, neighbourItem.Value.d}]; visited[pointv{point, neighbourItem.Value.d}] && !exists {
				continue
			}

			if vector, exists := m.allBestPaths[pointv{point, neighbourItem.Value.d}]; exists {
				if vector.totalCost == totalCost {
					m.storeAllBestPath(neighbourItem.Value)
				}
			}

			heap.Push(&pq, neighbourItem)

		}
	}
	return -1
}
