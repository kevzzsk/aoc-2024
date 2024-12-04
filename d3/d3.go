package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func readInput() string {
	// read from txt file
	file, err := os.Open("input3.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var data []string

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		data = append(data, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return strings.Join(data, "")
}

func main() {
	startTime := time.Now()
	stringToParse := readInput()

	p1(stringToParse)
	p2(stringToParse)

	endTime := time.Since(startTime)
	fmt.Printf("Day3 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds())
}

func p1(data string) {
	// match mul(1,2) and extract 1 and 2
	re := regexp.MustCompile(`(?:mul\()(\d+),(\d+)(?:\))`)

	pairs := re.FindAllStringSubmatch(data, -1)

	total := 0
	for _, match := range pairs {
		// e.g. match: [mul(1,2) 1 2]
		operand1, err := strconv.Atoi(match[1])
		if err != nil {
			panic(err)
		}
		operand2, err := strconv.Atoi(match[2])
		if err != nil {
			panic(err)
		}
		total += operand1 * operand2
	}

	fmt.Println("Total: ", total)
}

func p2(data string) {
	// match mul(1,2) and do() and don't()
	re := regexp.MustCompile(`((?:mul\()\d+,\d+(?:\))|(do\(\))|(don't\(\)))`)

	// match mul(1,2) and extract 1 and 2 - same as part 1
	getOperands := regexp.MustCompile(`(?:mul\()(\d+),(\d+)(?:\))`)

	pairs := re.FindAllString(data, -1)

	total := 0
	isEnabled := true // starts with enabled

	for _, match := range pairs {
		// e.g. match: [mul(1,2) do() don't() mul(3,4)]
		if match == "do()" {
			isEnabled = true
		} else if match == "don't()" {
			isEnabled = false
		} else {
			if isEnabled {
				// guaranteed to have 2 operands as match contains `mul(X,Y)`
				operands := getOperands.FindStringSubmatch(match)
				// e.g. operands: [mul(1,2) 1 2]
				operand1, err := strconv.Atoi(operands[1])
				if err != nil {
					panic(err)
				}
				operand2, err := strconv.Atoi(operands[2])
				if err != nil {
					panic(err)
				}
				total += operand1 * operand2
			}
		}
	}

	fmt.Println("Total: ", total)
}
