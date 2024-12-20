package main

import (
	utils "aoc-2024"
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type point struct {
	x int
	y int
}

type cheatPoint struct {
	start point
	end   point
}

type maze struct {
	grid        [][]string
	startPoint  point
	endPoint    point
	bestPath    map[point]Vector
	bestPathArr []Vector
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
	file, err := os.Open("input20.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	maze := maze{
		bestPath: make(map[point]Vector),
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

	for iy, row := range m.grid {
		for ix, v := range row {
			if ix == m.startPoint.x && iy == m.startPoint.y {
				fmt.Print("S")
			} else if ix == m.endPoint.x && iy == m.endPoint.y {
				fmt.Print("E")
			} else if _, exists := m.bestPath[point{ix, iy}]; exists {
				fmt.Print("X")
			} else {
				fmt.Print(v)
			}
		}
		fmt.Println()
	}
}

func main() {
	startTime := time.Now()

	maze := readInput()

	maze.traverse()

	// part 1 - cheat picoseconds is 2
	cheatsCount := make(map[cheatPoint]bool, 0)
	for _, v := range maze.bestPath {
		// for each best path check if 100 steps ahead there is any point within 2 steps
		cheatablePoint := maze.getCheatablePoint(v, 2)
		for _, cp := range cheatablePoint {
			cheatsCount[cp] = true
		}
	}

	fmt.Println("Cheats count: ", len(cheatsCount))

	// part 2 - cheat picoseconds is 20
	cheatsCount2 := make(map[cheatPoint]bool, 0)
	for _, v := range maze.bestPath {
		// for each best path check if 100 steps ahead there is any point within 2 steps
		cheatablePoint := maze.getCheatablePoint(v, 20)
		for _, cp := range cheatablePoint {
			cheatsCount2[cp] = true
		}
	}

	fmt.Println("Cheats count 2: ", len(cheatsCount2))

	endTime := time.Since(startTime)
	fmt.Printf("Day20 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day20 execution took: 1122 ms (1122870 µs)
}

func (m *maze) getStepsTaken(p1 point, p2 point) int {
	return utils.Abs(p1.x-p2.x) + utils.Abs(p1.y-p2.y)
}

// get all bestpath points that are n steps away from the current point
func (m *maze) getCheatablePoint(v Vector, cheatSeconds int) []cheatPoint {
	cheatablePoints := []cheatPoint{}
	const MIN_SAVED_PICOSECONDS = 100
	currentStep := v.step
	for i := len(m.bestPathArr) - 1; i >= currentStep+MIN_SAVED_PICOSECONDS; i-- {
		stepsTaken := m.getStepsTaken(v.point, m.bestPathArr[i].point)
		if stepsTaken <= cheatSeconds {
			// check actual picoseconds saved
			picosecondsSaved := (m.bestPathArr[i].step - currentStep) - stepsTaken
			if picosecondsSaved >= MIN_SAVED_PICOSECONDS {
				cheatablePoints = append(cheatablePoints, cheatPoint{start: v.point, end: m.bestPathArr[i].point})

			}
		}

	}
	return cheatablePoints
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
			m.storeBestPath(currentVector)
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

func (m *maze) storeBestPath(currentVector Vector) {
	currentVectorPtr := &currentVector
	for currentVectorPtr != nil {
		m.bestPath[currentVectorPtr.point] = *currentVectorPtr
		m.bestPathArr = append([]Vector{*currentVectorPtr}, m.bestPathArr...)
		currentVectorPtr = currentVectorPtr.parent
	}
}
