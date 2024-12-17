package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type warehouse struct {
	robotX int
	robotY int
	grid   [][]string
	moves  []string
}

// this function reads the input from the txt file and splits the rules and prints
func readInput(twiceAsWide bool) warehouse {
	// read from txt file
	file, err := os.Open("input15.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var warehouse warehouse

	// read the file line by line
	scanner := bufio.NewScanner(file)
	readingMap := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingMap = false
			continue
		}
		row := strings.Split(line, "")
		if readingMap {
			currRow := make([]string, 0)
			for i, item := range row {
				if twiceAsWide {
					if item == "#" || item == "." {
						currRow = append(currRow, item)
						currRow = append(currRow, item)
					} else if item == "O" {
						currRow = append(currRow, "[")
						currRow = append(currRow, "]")
					} else if item == "@" {
						currRow = append(currRow, "@")
						currRow = append(currRow, ".")
					}
				} else {
					currRow = append(currRow, item)
				}

				if item == "@" {
					if twiceAsWide {
						warehouse.robotX = i * 2
						warehouse.robotY = len(warehouse.grid)
					} else {
						warehouse.robotX = i
						warehouse.robotY = len(warehouse.grid)
					}
				}
			}
			warehouse.grid = append(warehouse.grid, currRow)
		} else {
			warehouse.moves = append(warehouse.moves, row...)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return warehouse
}

func main() {
	startTime := time.Now()
	w := readInput(false)

	fmt.Printf("%+v\n", w)

	for _, move := range w.moves {
		w.moveRobot(move)
	}
	fmt.Println("warehouse sum of GPS coordinates:", w.findGPSCoord())

	w2 := readInput(true)
	for _, move := range w2.moves {
		w2.moveRobotV2(move)
	}
	w2.printWarehouse()
	fmt.Println("warehouse 2 sum of GPS coordinates:", w2.findGPSCoordV2())

	endTime := time.Since(startTime)
	fmt.Printf("Day15 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day15 execution took: 40 ms (40626 µs)
}

func (w *warehouse) printWarehouse() {
	for _, row := range w.grid {
		for _, cell := range row {
			fmt.Print(cell)
		}
		fmt.Println()
	}
}

func (w *warehouse) moveRobot(direction string) {
	freeX, freeY := w.getNextAvailableSpot(direction)
	if freeX != -1 && freeY != -1 {
		// move robot + items its pushing
		switch direction {
		case "^":
			temp := w.grid[freeY][w.robotX]
			for y := freeY; y != w.robotY; y++ {
				w.grid[y][w.robotX] = w.grid[y+1][w.robotX]
			}
			w.grid[w.robotY][w.robotX] = temp
			w.robotY -= 1
		case ">":
			temp := w.grid[w.robotY][freeX]
			for x := freeX; x != w.robotX; x-- {
				w.grid[w.robotY][x] = w.grid[w.robotY][x-1]
			}
			w.grid[w.robotY][w.robotX] = temp
			w.robotX += 1
		case "v":
			temp := w.grid[freeY][w.robotX]
			for y := freeY; y != w.robotY; y-- {
				w.grid[y][w.robotX] = w.grid[y-1][w.robotX]
			}
			w.grid[w.robotY][w.robotX] = temp
			w.robotY += 1
		case "<":
			temp := w.grid[w.robotY][freeX]
			for x := freeX; x != w.robotX; x++ {
				w.grid[w.robotY][x] = w.grid[w.robotY][x+1]
			}
			w.grid[w.robotY][w.robotX] = temp
			w.robotX -= 1

		}
	}
}

func (w *warehouse) getNextAvailableSpot(direction string) (int, int) {
	x := w.robotX
	y := w.robotY
	switch direction {
	case "^":
		for ; y > 0; y-- {
			if w.grid[y][x] == "." {
				return x, y
			}
			if w.grid[y][x] == "#" {
				return -1, -1
			}
		}
	case ">":
		for ; x < len(w.grid[y]); x++ {
			if w.grid[y][x] == "." {
				return x, y
			}
			if w.grid[y][x] == "#" {
				return -1, -1
			}
		}
	case "v":
		for ; y < len(w.grid); y++ {
			if w.grid[y][x] == "." {
				return x, y
			}
			if w.grid[y][x] == "#" {
				return -1, -1
			}
		}
	case "<":
		for ; x > 0; x-- {
			if w.grid[y][x] == "." {
				return x, y
			}
			if w.grid[y][x] == "#" {
				return -1, -1
			}
		}
	}
	return -1, -1
}

func (w *warehouse) findGPSCoord() int {
	sum := 0
	for y, row := range w.grid {
		for x, cell := range row {
			if cell == "O" {
				sum += (y * 100) + x
			}
		}
	}
	return sum
}

func (w *warehouse) findGPSCoordV2() int {
	sum := 0
	for y, row := range w.grid {
		for x, cell := range row {
			if cell == "[" {
				sum += (y * 100) + x
			}
		}
	}
	return sum
}

type node struct {
	x       int
	y       int
	element string
	front   *node
	back    *node
	sibling *node
}

type tree struct {
	root *node
}

func (t *tree) printTree() {
	if t.root == nil {
		return
	}
	queue := make([]*node, 0)
	queue = append(queue, t.root)
	for len(queue) > 0 {
		currNode := queue[0]
		queue = queue[1:]
		fmt.Printf("x: %v, y: %v, element: %v\n", currNode.x, currNode.y, currNode.element)
		if currNode.front != nil {
			queue = append(queue, currNode.front)
		}
		if currNode.sibling != nil {
			queue = append(queue, currNode.sibling)
		}
	}
}

// this fn creates a tree of the items connected to the robot in a BFS manner.
// if all leaf nodes are "." then we have space to push item (if any) and move forward.
// if any leaf node is not "." then we can't move forward.
// moving forward just means moving item from the tree nodes 1 level down.
// we start from the leaf nodes and move the item to its front node.
// for cases of ">" and "<" it is a linear move as only 1 row is involved.
// for cases of "^" and "v" it will be messy as boxes can chain in a massive configuration
func (w *warehouse) moveRobotV2(direction string) tree {
	// create root node
	root := &node{
		x:       w.robotX,
		y:       w.robotY,
		element: w.grid[w.robotY][w.robotX],
	}
	// create tree
	t := tree{root: root}

	leafNodes := make([]*node, 0)

	queue := make([]*node, 0)
	queue = append(queue, root)
	// traverse tree and add nodes
	for len(queue) > 0 {
		currNode := queue[0]
		currX := currNode.x
		currY := currNode.y
		queue = queue[1:]
		if direction == "^" {
			frontElement := w.grid[currY-1][currX]
			if frontElement == "]" {
				currNode.front = &node{
					x:       currX,
					y:       currY - 1,
					element: "]",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
				currNode.front.sibling = &node{
					x:       currX - 1,
					y:       currY - 1,
					element: "[",
				}
				queue = append(queue, currNode.front.sibling)
			} else if frontElement == "[" {
				currNode.front = &node{
					x:       currX,
					y:       currY - 1,
					element: "[",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
				currNode.front.sibling = &node{
					x:       currX + 1,
					y:       currY - 1,
					element: "]",
				}
				queue = append(queue, currNode.front.sibling)
			} else if frontElement == "." || frontElement == "#" {
				currNode.front = &node{
					x:       currX,
					y:       currY - 1,
					element: frontElement,
					back:    currNode,
				}
				leafNodes = append(leafNodes, currNode.front)
			}
		}
		if direction == ">" {
			frontElement := w.grid[currY][currX+1]
			if frontElement == "]" {
				currNode.front = &node{
					x:       currX + 1,
					y:       currY,
					element: "]",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
			} else if frontElement == "[" {
				currNode.front = &node{
					x:       currX + 1,
					y:       currY,
					element: "[",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
			} else if frontElement == "." || frontElement == "#" {
				currNode.front = &node{
					x:       currX + 1,
					y:       currY,
					element: frontElement,
					back:    currNode,
				}
				leafNodes = append(leafNodes, currNode.front)
			}
		}
		if direction == "v" {
			frontElement := w.grid[currY+1][currX]
			if frontElement == "]" {
				currNode.front = &node{
					x:       currX,
					y:       currY + 1,
					element: "]",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
				currNode.front.sibling = &node{
					x:       currX - 1,
					y:       currY + 1,
					element: "[",
				}
				queue = append(queue, currNode.front.sibling)
			} else if frontElement == "[" {
				currNode.front = &node{
					x:       currX,
					y:       currY + 1,
					element: "[",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
				currNode.front.sibling = &node{
					x:       currX + 1,
					y:       currY + 1,
					element: "]",
				}
				queue = append(queue, currNode.front.sibling)
			} else if frontElement == "." || frontElement == "#" {
				currNode.front = &node{
					x:       currX,
					y:       currY + 1,
					element: frontElement,
					back:    currNode,
				}
				leafNodes = append(leafNodes, currNode.front)
			}
		}
		if direction == "<" {
			frontElement := w.grid[currY][currX-1]
			if frontElement == "]" {
				currNode.front = &node{
					x:       currX - 1,
					y:       currY,
					element: "]",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
			} else if frontElement == "[" {
				currNode.front = &node{
					x:       currX - 1,
					y:       currY,
					element: "[",
					back:    currNode,
				}
				queue = append(queue, currNode.front)
			} else if frontElement == "." || frontElement == "#" {
				currNode.front = &node{
					x:       currX - 1,
					y:       currY,
					element: frontElement,
					back:    currNode,
				}
				leafNodes = append(leafNodes, currNode.front)
			}
		}
	}

	// if all leaf nodes are "." then we can move forward
	canTraverse := true
	for _, leaf := range leafNodes {
		if leaf.element != "." {
			canTraverse = false
		}
	}

	backQueue := make([]*node, 0)
	backQueue = append(backQueue, leafNodes...)
	if canTraverse {
		// move all nodes to its front node
		for len(backQueue) > 0 {
			currNode := backQueue[0]
			backQueue = backQueue[1:]
			if currNode.back != nil {
				// move element to front
				w.grid[currNode.y][currNode.x] = currNode.back.element
				// set front element to empty space
				w.grid[currNode.back.y][currNode.back.x] = "."
				backQueue = append(backQueue, currNode.back)
			}
		}
		// update robot location
		w.robotX = root.front.x
		w.robotY = root.front.y
		// set the robot location to empty space
		w.grid[root.y][root.x] = "."
	}
	return t
}
