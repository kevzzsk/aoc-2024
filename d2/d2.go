package main

import (
	utils "aoc-2024"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readInput() [][]int {
	// read from txt file
	file, err := os.Open("input2.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var reports [][]int

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		report := []int{}
		for _, part := range parts {
			num, err := strconv.Atoi(part)
			if err != nil {
				panic(err)
			}
			report = append(report, num)
		}
		reports = append(reports, report)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return reports
}

func main() {
	startTime := time.Now()
	reports := readInput()

	p1(reports)
	p2(reports)

	endTime := time.Since(startTime)
	fmt.Printf("Day2 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds())
}

// report is safe when all the levels (digits) are
// - all either increasing or all decreasing
// - two adjacent digits differ by at least 1 and at most 3
func isSafeReport(report []int, problemDampenerCount int) bool {
	isIncreasing := report[0] < report[1]

	for i := 1; i < len(report); i++ {
		diff := report[i-1] - report[i]
		if (isIncreasing && diff > 0) || (!isIncreasing && diff < 0) {
			// check 1st criteria
			if problemDampenerCount == 0 {
				return false
			}

			// In special scenario of a decrease-increase or increase-decrease case, there are 3 levels (digits) being affected
			// i-2, i-1 and i
			// This special scenario can only happen in the first 3 levels of a report (i.e [20 21 19 ...])
			// in this case, we will only encounter issue when checking for i-1 and i level. Hence we must backtrack i-2 as well

			// try removing i-1th item
			isSafeRemove1 := isSafeReport(utils.RemoveElement(report, i-1), problemDampenerCount-1)
			// try removing ith item
			isSafeRemove2 := isSafeReport(utils.RemoveElement(report, i), problemDampenerCount-1)
			// try removing i-2th item
			isSafeRemove3 := isSafeReport(utils.RemoveElement(report, i-2), problemDampenerCount-1)
			return isSafeRemove1 || isSafeRemove2 || isSafeRemove3

		} else if utils.Abs(diff) > 3 || utils.Abs(diff) < 1 {
			// check 2nd criteria
			if problemDampenerCount == 0 {
				return false
			}

			// try removing i-1th item
			isSafeRemove1 := isSafeReport(utils.RemoveElement(report, i-1), problemDampenerCount-1)
			// try removing ith item
			isSafeRemove2 := isSafeReport(utils.RemoveElement(report, i), problemDampenerCount-1)
			return isSafeRemove1 || isSafeRemove2
		}
	}
	return true
}

// find how many safe reports there are
// time complexity: O(n*m)
// space complexity: O(1)
func p1(reports [][]int) {
	safeCount := 0

	for _, report := range reports {
		if isSafe := isSafeReport(report, 0); isSafe {
			safeCount++
		}
	}

	fmt.Println("Safe Reports: ", safeCount)
}

// find how many safe reports there are
// time complexity: O(n*m)
// space complexity: O(1)
func p2(reports [][]int) {
	safeCount := 0

	for _, report := range reports {
		if isSafe := isSafeReport(report, 1); isSafe {
			safeCount++
		}
	}

	fmt.Println("Safe Reports with problem dampener: ", safeCount)
}
