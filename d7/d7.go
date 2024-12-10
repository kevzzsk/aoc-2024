package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type EquationInput struct {
	result   int
	operands []int
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() []EquationInput {
	// read from txt file
	file, err := os.Open("input7.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var equations []EquationInput

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Split(line, ":")
		operands := strings.Split(strings.TrimSpace(row[1]), " ")
		// convert the operands to int
		var operandsInt []int
		for _, operand := range operands {
			operandInt, _ := strconv.Atoi(operand)
			operandsInt = append(operandsInt, operandInt)
		}
		resultInt, _ := strconv.Atoi(row[0])
		equations = append(equations, EquationInput{resultInt, operandsInt})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return equations
}

func main() {
	startTime := time.Now()
	equations := readInput()

	// PART 1
	// find the number of equations that are valid
	sumOfValidEquations := 0
	for _, equation := range equations {
		if backtracking(equation, false) {
			sumOfValidEquations += equation.result
		}
	}

	fmt.Println("Sum of valid equations: ", sumOfValidEquations)

	// PART 2
	// find the number of equations that are valid
	sumOfValidEquationsWithConcat := 0
	for _, equation := range equations {
		if backtracking(equation, true) {
			sumOfValidEquationsWithConcat += equation.result
		}
	}

	fmt.Println("Sum of valid equations with concat: ", sumOfValidEquationsWithConcat)

	endTime := time.Since(startTime)
	fmt.Printf("Day7 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) //Day7 execution took: 2 ms (2123 µs)
}

func isDivisible(a, b int) bool {
	if a == 0 || b == 0 {
		return false
	}
	return a%b == 0
}

// find if last digits of A is B
func isConcatenation(a, b int) bool {
	if a == 0 || b == 0 {
		return false
	}
	for b > 0 {
		if a%10 != b%10 {
			return false
		}
		a /= 10
		b /= 10
	}
	return true
}

// Check from the last element of the operands and try if it is possible to reach the result
// always check for division first, then concatenation and then subtraction
// if the last operand is divisible by the result, then divide it
// if the last operand is concatenation of the result, then de-concatenate it
// if the last operand is not divisible by the result, then subtract it
// if none of the above is possible, then return false
func backtracking(equation EquationInput, useConcatenation bool) bool {
	//base case
	if equation.result < 0 {
		return false
	}
	if len(equation.operands) == 1 {
		return equation.operands[0] == equation.result
	}

	lastOperand := equation.operands[len(equation.operands)-1]
	if isDivisible(equation.result, lastOperand) {
		if backtracking(EquationInput{equation.result / lastOperand, equation.operands[:len(equation.operands)-1]}, useConcatenation) {
			return true
		}
	}
	if useConcatenation && isConcatenation(equation.result, lastOperand) {
		newResult := equation.result
		newLastOperange := lastOperand
		for newLastOperange > 0 {
			newResult /= 10
			newLastOperange /= 10
		}
		if backtracking(EquationInput{newResult, equation.operands[:len(equation.operands)-1]}, useConcatenation) {
			return true
		}
	}
	// use "+" operator
	if backtracking(EquationInput{equation.result - lastOperand, equation.operands[:len(equation.operands)-1]}, useConcatenation) {
		return true
	}
	return false
}
