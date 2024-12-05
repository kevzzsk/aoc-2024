package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

// this function reads the input from the txt file and splits the rules and prints
func readInput() ([][]string, [][]string) {
	// read from txt file
	file, err := os.Open("input5.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var rules [][]string
	var prints [][]string

	readingRules := true
	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingRules = false
			continue
		}
		if readingRules {
			rules = append(rules, strings.Split(line, "|"))
		} else {
			prints = append(prints, strings.Split(line, ","))
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return rules, prints
}

func main() {
	startTime := time.Now()
	rules, prints := readInput()

	var invalidPrints [][]string

	p1(rules, prints, &invalidPrints)
	p2(rules, &invalidPrints)

	endTime := time.Since(startTime)
	fmt.Printf("Day5 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day5 execution took: 8 ms (8752 µs)
}

// Find correctly ordered prints from a set of rules and total up the middle number of the array
//
// Approached used is to create a ruleset containing the key and the values that can follow it (to its right)
// Then for each print, check if the next item (i) is in the ruleset of the previous item (i-1)
// If all items in the print are valid, add the middle number to the total sum of valid prints
// Don't have to check each item against all the rest, just check if the next item is in the ruleset of the previous item
// e.g. [A B C D E], just check if B is in the ruleset of A, then C in the ruleset of B, and so on
// This works because the order is transitive (if A -> B and B -> C, then A -> C)
func p1(rules [][]string, prints [][]string, invalidPrints *[][]string) {
	ruleSet := make(map[string]map[string]bool)

	// init ruleSet
	for _, rule := range rules {
		if _, exists := ruleSet[rule[0]]; exists {
			ruleSet[rule[0]][rule[1]] = true
		} else {
			ruleSet[rule[0]] = make(map[string]bool)
			ruleSet[rule[0]][rule[1]] = true
		}
	}

	totalSumOfValidPrints := 0
	// check prints against ruleset
	for _, print := range prints {
		midOfPrintString := print[len(print)/2]
		midOfPrint, err := strconv.Atoi(midOfPrintString)
		if err != nil {
			panic(err)
		}
		for i := 1; i < len(print); i++ {
			// get i-1 rule
			rule := ruleSet[print[i-1]]
			// check if i exists in i-1 rule
			if _, exists := rule[print[i]]; !exists {
				// add invalid prints to slice for part2
				*invalidPrints = append(*invalidPrints, print)
				break
			}
			if i == len(print)-1 {
				totalSumOfValidPrints += midOfPrint
			}
		}
	}

	fmt.Println("Total sum of valid prints: ", totalSumOfValidPrints)
}

// Fix the invalid prints by sorting them in the correct order and summing up the middle number of the array
//
// Approach used is to sort the invalid prints by counting how many items goes to its right side
func p2(rules [][]string, invalidPrints *[][]string) {
	ruleSet := make(map[string]map[string]bool)

	// init ruleSet
	for _, rule := range rules {
		if _, exists := ruleSet[rule[0]]; exists {
			ruleSet[rule[0]][rule[1]] = true
		} else {
			ruleSet[rule[0]] = make(map[string]bool)
			ruleSet[rule[0]][rule[1]] = true
		}
	}

	totalSumOfNewValidPrints := 0
	// sort number by counting how many numbers goes to its right side
	// e.g. 10|11 10|12 11|12
	// [12 10 11]
	// for 10, how many (11,12) it has on its ruleset -> 2 (because both 11 and 12 are present)
	// for 11, -> 1 (only 12)
	// for 12 -> 0 (because 12 doesnt not appear on the leftside)
	// hence the correct order is [10 11 12]
	for _, print := range *invalidPrints {

		// compare the count of items that goes before it
		comp := func(a, b string) int {
			countA := countItemsInRuleSet(a, ruleSet, print)
			countB := countItemsInRuleSet(b, ruleSet, print)
			if countA < countB {
				return -1
			}
			return 1
		}

		slices.SortFunc(print, comp)

		midIndex := len(print) / 2
		intVal, _ := strconv.Atoi(print[midIndex])
		totalSumOfNewValidPrints += intVal
	}

	fmt.Println("Total sum of new valid prints: ", totalSumOfNewValidPrints)
}

func countItemsInRuleSet(key string, ruleSet map[string]map[string]bool, print []string) int {
	count := 0
	for _, item := range print {
		if ruleSet[key][item] {
			count++
		}
	}
	return count
}
