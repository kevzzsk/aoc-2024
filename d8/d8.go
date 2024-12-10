package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Grid struct {
	grid                       [][]string
	antennas                   map[string][]point
	antiNodesLocations         map[point]bool
	resonantAntiNodesLocations map[point]bool
}

func (g *Grid) isPointInGrid(p point) bool {
	return p.x >= 0 && p.x < len(g.grid[0]) && p.y >= 0 && p.y < len(g.grid)
}

func getAntiNodes(p1, p2 point) (point, point) {
	xDiff := p1.x - p2.x
	yDiff := p1.y - p2.y
	return point{p1.x + xDiff, p1.y + yDiff}, point{p2.x + (xDiff * -1), p2.y + (yDiff * -1)}

}

func (g *Grid) findAllAntiNodes() {

	for _, antenna := range g.antennas {
		for ip, p := range antenna {
			for _, p2 := range antenna[ip+1:] {
				antiNode1, antiNode2 := getAntiNodes(p, p2)
				if g.isPointInGrid(antiNode1) {
					g.antiNodesLocations[antiNode1] = true
				}
				if g.isPointInGrid(antiNode2) {
					g.antiNodesLocations[antiNode2] = true
				}
			}
		}
	}
}

func (g *Grid) getResonantAntiNodes(p1, p2 point) []point {

	resonantAntiNodes := make([]point, 0)
	xDiff := p1.x - p2.x
	yDiff := p1.y - p2.y

	// include the two points
	resonantAntiNodes = append(resonantAntiNodes, p1, p2)

	// generate all resonant anti nodes that are in the grid
	for i := 1; g.isPointInGrid(point{p1.x + (xDiff * i), p1.y + (yDiff * i)}); i++ {
		resonantAntiNodes = append(resonantAntiNodes, point{p1.x + (xDiff * i), p1.y + (yDiff * i)})
	}
	// generate in the opposite direction
	for i := 1; g.isPointInGrid(point{p2.x + (xDiff * -1 * i), p2.y + (yDiff * -1 * i)}); i++ {
		resonantAntiNodes = append(resonantAntiNodes, point{p2.x + (xDiff * -1 * i), p2.y + (yDiff * -1 * i)})
	}
	return resonantAntiNodes
}

func (g *Grid) findAllResonantAntiNodes() {

	for _, antenna := range g.antennas {
		for ip, p := range antenna {
			for _, p2 := range antenna[ip+1:] {
				resonantAntiNodes := g.getResonantAntiNodes(p, p2)
				for _, resonantAntiNode := range resonantAntiNodes {
					g.resonantAntiNodesLocations[resonantAntiNode] = true
				}
			}
		}
	}
}

func (g *Grid) printGrid() {
	for iy, row := range g.grid {
		for ix, item := range row {
			if item != "." {
				fmt.Print(item)
			} else if _, ok := g.antiNodesLocations[point{ix, iy}]; ok {
				fmt.Print("#")
			} else if _, ok := g.resonantAntiNodesLocations[point{ix, iy}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(item)
			}
		}
		fmt.Println()
	}
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() [][]string {
	// read from txt file
	file, err := os.Open("input8.txt")
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

type point struct {
	x int
	y int
}

func main() {
	startTime := time.Now()
	grid := readInput()

	// find all antennas locations
	locations := make(map[string][]point)
	for iy, row := range grid {
		for ix, item := range row {
			if item != "." {
				locations[item] = append(locations[item], point{ix, iy})
			}
		}
	}

	g := Grid{grid, locations, make(map[point]bool), make(map[point]bool)}
	g.findAllAntiNodes()
	g.findAllResonantAntiNodes()
	// g.printGrid()

	fmt.Println("Antinodes Unique Locations: ", len(g.antiNodesLocations))
	fmt.Println("Resonant Antinodes Unique Locations: ", len(g.resonantAntiNodesLocations))

	endTime := time.Since(startTime)
	fmt.Printf("Day8 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) //Day8 execution took: 1 ms (1179 µs)
}
