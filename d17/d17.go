package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type device struct {
	registerA          int
	registerB          int
	registerC          int
	program            []int
	instructionPointer int
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() device {
	// read from txt file
	file, err := os.Open("input17.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	device := device{
		instructionPointer: 0,
	}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			index++
			continue
		}
		if index == 0 {
			registerA, _ := strconv.Atoi(strings.Split(line, " ")[2])
			device.registerA = registerA
		} else if index == 1 {
			registerB, _ := strconv.Atoi(strings.Split(line, " ")[2])
			device.registerB = registerB
		} else if index == 2 {
			registerC, _ := strconv.Atoi(strings.Split(line, " ")[2])
			device.registerC = registerC
		} else {
			program := strings.Split(strings.Split(line, " ")[1], ",")
			for _, p := range program {
				pInt, _ := strconv.Atoi(p)
				device.program = append(device.program, pInt)
			}
		}
		index++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return device
}

func main() {
	startTime := time.Now()
	d := readInput()

	// part 1 - pretty straight forward
	fmt.Println(d.run())

	// part 2 find initial registerA value that produces the same program output
	programStr := make([]string, len(d.program))
	for i, v := range d.program {
		programStr[i] = strconv.Itoa(v)
	}
	programStrStr := strings.Join(programStr, ",")

	// for part 2, alot of manual work to find the bounding range as only certain number range produces 16 digits outputs.
	// essentially, bigger registerA will produce more outputs.
	// using an exponential increase to find the first instance of 16 digits output (lower bound) - in this case its 35_184_372_088_832
	// then find the next instance of 17 digits output (upper bound) - in this case its 281_474_976_710_656
	// then iterate with exponential increase based on numbers of matches from the end of the output to the end of the program
	// this is because the end of the output covers a larger range of possible registerA values
	for currentRegisterA := 35_184_372_088_832; currentRegisterA < 281_474_976_710_656; currentRegisterA++ {

		dcopy := device{
			registerA:          currentRegisterA,
			registerB:          d.registerB,
			registerC:          d.registerC,
			program:            d.program,
			instructionPointer: 0,
		}
		res := dcopy.run()
		arrRes := strings.Split(res, ",")
		matches := reverseArrMatchesCount(programStr, arrRes)

		// for debugging
		if matches > 6 {
			fmt.Println("matches:", matches, "registerA:", formatNumber(currentRegisterA), converDecToBin(currentRegisterA), res)
		}
		if res == programStrStr {
			fmt.Println("matches:", matches, "registerA:", formatNumber(currentRegisterA), res)
			break
		}
		// for exponential increase based on matches
		currentRegisterA += int(math.Pow(2.5, float64(16-matches)))
	}

	endTime := time.Since(startTime)
	fmt.Printf("Day17 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds())
}

func formatNumber(n int) string {
	// Convert number to string
	str := strconv.Itoa(n)

	// Handle numbers less than 1000
	if len(str) <= 3 {
		return str
	}

	// Insert commas
	var result []byte
	for i, j := len(str)-1, 0; i >= 0; i-- {
		if j > 0 && j%3 == 0 {
			result = append([]byte{'_'}, result...)
		}
		result = append([]byte{str[i]}, result...)
		j++
	}

	return string(result)
}

func converDecToOctal(n int) string {
	return strconv.FormatInt(int64(n), 8)
}

func converDecToBin(n int) string {
	return strconv.FormatInt(int64(n), 2)
}

// check how many elements from end of arr1 and end of arr2 matches
// works for arrays of different lengths
func reverseArrMatchesCount(arr1, arr2 []string) int {
	matches := 0
	for i := 0; i < len(arr1) && i < len(arr2); i++ {
		if arr1[len(arr1)-1-i] == arr2[len(arr2)-1-i] {
			matches++
		} else {
			break
		}
	}
	return matches
}

func (d *device) getComboOperand(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return d.registerA
	case 5:
		return d.registerB
	case 6:
		return d.registerC
	default:
		return operand
	}
}

func (d *device) getOpcodeAndOperand() (opcode, operand int) {
	if d.instructionPointer >= len(d.program) {
		return -1, -1
	}
	opcode = d.program[d.instructionPointer]
	operand = d.program[d.instructionPointer+1]
	return opcode, operand
}

func (d *device) run() string {
	output := []string{}
	for opcode, operand := d.getOpcodeAndOperand(); opcode != -1; opcode, operand = d.getOpcodeAndOperand() {
		switch opcode {
		case 0:
			// adv
			d.registerA = d.registerA / int(math.Pow(2, float64(d.getComboOperand(operand))))
		case 1:
			//bxl
			d.registerB = d.registerB ^ operand
		case 2:
			// bst
			d.registerB = d.getComboOperand(operand) % 8
		case 3:
			// jnz
			if d.registerA != 0 {
				d.instructionPointer = operand
				continue
			}
		case 4:
			//bxc
			d.registerB = d.registerB ^ d.registerC
		case 5:
			//out
			res := d.getComboOperand(operand) % 8
			output = append(output, strconv.Itoa(res))
		case 6:
			//bdv
			d.registerB = d.registerA / int(math.Pow(2, float64(d.getComboOperand(operand))))
		case 7:
			//cdv
			d.registerC = d.registerA / int(math.Pow(2, float64(d.getComboOperand(operand))))
		}
		d.instructionPointer += 2
	}
	return strings.Join(output, ",")
}
