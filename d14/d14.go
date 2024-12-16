package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type robot struct {
	x, y   int
	vx, vy int
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() []robot {
	// read from txt file
	file, err := os.Open("input14.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var robots []robot

	re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)
		if len(matches) == 5 {
			x, _ := strconv.Atoi(matches[1])
			y, _ := strconv.Atoi(matches[2])
			vx, _ := strconv.Atoi(matches[3])
			vy, _ := strconv.Atoi(matches[4])
			robots = append(robots, robot{x: x, y: y, vx: vx, vy: vy})
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return robots
}

func main() {
	startTime := time.Now()
	robots := readInput()

	quadrantCount := make(map[int]int)

	for _, robot := range robots {
		robot.runForNSeconds(100)
		if robot.getQuadrant() != -1 {
			quadrantCount[robot.getQuadrant()]++
		}
	}

	// multiple all the quadrants count together
	safetyFactor := 1
	for _, count := range quadrantCount {
		safetyFactor *= count
	}
	fmt.Println("quadrantCount:", quadrantCount)
	fmt.Println("Safety Factor:", safetyFactor)

	// reset
	robots = readInput()

	// part 2 find christmas tree
	robotPositions := make(map[string]bool)
	seconds := 0
	for {
		robotPositions = make(map[string]bool)
		for i := range robots {
			robots[i].runForNSeconds(1)
			robotPositions[fmt.Sprintf("%d,%d", robots[i].x, robots[i].y)] = true
		}
		// assumption: all the robots must not be overlapping for it to form the christmas tree
		if len(robotPositions) == len(robots) {
			fmt.Println("len(robotPositions):", len(robotPositions), "len(robots):", len(robots), "seconds:", seconds)
			printMap(robotPositions)
			break
		}
		seconds += 1

	}

	endTime := time.Since(startTime)
	fmt.Printf("Day14 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day14 execution took: 642 ms (642481 µs)
}

func printMap(robotPositions map[string]bool) {
	// print a grid of the robots
	WIDTH, HEIGHT := getMapSize()
	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			if robotPositions[fmt.Sprintf("%d,%d", x, y)] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func getMapSize() (int, int) {
	const WIDTH = 101
	const HEIGHT = 103

	return WIDTH, HEIGHT
}

func (r *robot) runForNSeconds(n int) {
	WIDTH, HEIGHT := getMapSize()

	x := (r.x + r.vx*n) % WIDTH
	y := (r.y + r.vy*n) % HEIGHT

	if x < 0 {
		x += WIDTH
	}
	if y < 0 {
		y += HEIGHT
	}

	// update the robot's position
	r.x = x
	r.y = y
}

// returns the quadrant the robot is in
// -1 for it being right on the border
// 0 for top left, 1 for top right, 2 for bottom left, 3 for bottom right
func (r *robot) getQuadrant() int {
	WIDTH, HEIGHT := getMapSize()

	// is in the top left
	if r.x < WIDTH/2 && r.y < HEIGHT/2 {
		return 0
	}
	// is in the top right
	if r.x > WIDTH/2 && r.y < HEIGHT/2 {
		return 1
	}
	// is in the bottom left
	if r.x < WIDTH/2 && r.y > HEIGHT/2 {
		return 2
	}
	// is in the bottom right
	if r.x > WIDTH/2 && r.y > HEIGHT/2 {
		return 3
	}
	return -1
}
