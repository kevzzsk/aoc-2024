package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// problem: https://adventofcode.com/2024/day/1

func readInput() ([]int, []int) {
	// read from txt file
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// create two lists to store the split values
	var list1, list2 []int

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		num1, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		num2, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		list1 = append(list1, num1)
		list2 = append(list2, num2)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return list1, list2
}

func absDiff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func main() {
	startTime := time.Now()

	list1, list2 := readInput()
	// sort the lists
	sort.Ints(list1)
	sort.Ints(list2)

	p1(list1, list2)
	p2(list1, list2)

	endTime := time.Since(startTime)
	fmt.Printf("Day1 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds())
}

// find the difference between the two lists of the same ranking from smallest to largest
// time complexity: O(n), since we sorted both list prior the real time complexity is O(nlogn)
// space complexity: O(1)
func p1(list1, list2 []int) {
	totalDiff := 0
	// find the difference between the two lists
	for i := 0; i < len(list1); i++ { // assumption is that the two lists are of the same length
		diff := absDiff(list1[i], list2[i])
		totalDiff += diff
	}
	fmt.Println("Total difference: ", totalDiff)
}

// find the similarity between the two lists
// time complexity: O(n + m), since we sorted both list prior the real time complexity is O(nlogn)
// space complexity: O(n)
func p2(list1, list2 []int) {

	similarityMultiplierMap := make(map[int]int)
	similarity := 0

	for _, val := range list2 {
		similarityMultiplierMap[val]++
	}

	for _, val := range list1 {
		similarity += similarityMultiplierMap[val] * val
	}

	fmt.Println("Similarity: ", similarity)
}
