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

type logic string

const (
	AND logic = "AND"
	OR  logic = "OR"
	XOR logic = "XOR"
)

type gate struct {
	inputA string
	inputB string
	output string
	logic  logic
}

type wires struct {
	state map[string]int
	gates []gate
}

func readInput() wires {
	// read from txt file
	file, err := os.Open("input24.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := wires{state: make(map[string]int)}

	// read the file line by line
	readingState := true
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingState = false
			continue
		}
		if readingState {
			row := strings.Split(line, ":")
			bitValue, _ := strconv.Atoi(strings.TrimSpace(row[1]))
			result.state[row[0]] = bitValue
		} else {
			row := strings.Split(line, " ")
			inputA := row[0]
			inputB := row[2]
			output := row[4]
			logic := logic(row[1])
			result.gates = append(result.gates, gate{inputA, inputB, output, logic})
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

func main() {
	startTime := time.Now()

	wires := readInput()

	bits := make([]int, 64)
	for _, gate := range wires.gates {
		if strings.HasPrefix(gate.output, "z") {
			// get the digits
			index := strings.SplitAfter(gate.output, "z")[1]
			bitIndex, _ := strconv.Atoi(index)
			bits[bitIndex] = wires.getOutput(gate.output)
		}
	}

	// part 2
	// 1. check all XOR gates, if both inputs are x or y, then output cannot be z except for x00 XOR y00 -> z00
	// 2. check all XOR gates, all other must output to z
	// 3. all gates that output to z must be XOR
	badGates := make([]gate, 0)
	for _, gate := range wires.gates {
		if gate.logic == XOR {
			if (strings.HasPrefix(gate.inputA, "x") || strings.HasPrefix(gate.inputA, "y")) && (strings.HasPrefix(gate.inputB, "x") || strings.HasPrefix(gate.inputB, "y")) {
				if strings.HasPrefix(gate.output, "z") && gate.output != "z00" {
					badGates = append(badGates, gate)
				}
			} else {
				if !strings.HasPrefix(gate.output, "z") {
					badGates = append(badGates, gate)
				}
			}

		}
		if strings.HasPrefix(gate.output, "z") && gate.logic != XOR && gate.output != "z45" {
			badGates = append(badGates, gate)
		}
		if gate.logic == AND && gate.inputA != "x00" && gate.inputB != "x00" {
			for _, subGate := range wires.gates {
				if (gate.output == subGate.inputA || gate.output == subGate.inputB) && subGate.logic != OR {
					badGates = append(badGates, gate)
				}
			}
		}

		if gate.logic == XOR {
			for _, subGate := range wires.gates {
				if (gate.output == subGate.inputA || gate.output == subGate.inputB) && subGate.logic == OR {
					badGates = append(badGates, gate)
				}
			}
		}
	}
	// remove dupes
	for i := 0; i < len(badGates); i++ {
		for j := i + 1; j < len(badGates); j++ {
			if badGates[i] == badGates[j] {
				badGates = append(badGates[:j], badGates[j+1:]...)
				j--
			}
		}
	}

	fmt.Println("badgates", badGates)

	// combine badgates output and sort
	badGatesOutput := make([]string, 0)
	for _, badGate := range badGates {
		badGatesOutput = append(badGatesOutput, badGate.output)
	}
	sort.Strings(badGatesOutput)
	fmt.Println("badgates output: ", strings.Join(badGatesOutput, ","))

	endTime := time.Since(startTime)
	fmt.Printf("Day24 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day24 execution took: 35 ms (35103 µs)
}

func (g gate) String() string {
	return fmt.Sprintf("%s %s %s = %s\n", g.inputA, g.logic, g.inputB, g.output)
}

func (w *wires) getOutput(wire string) int {
	if val, ok := w.state[wire]; ok {
		return val
	}

	for _, gate := range w.gates {
		if gate.output == wire {
			switch gate.logic {
			case AND:
				w.state[gate.output] = w.getOutput(gate.inputA) & w.getOutput(gate.inputB)
			case OR:
				w.state[gate.output] = w.getOutput(gate.inputA) | w.getOutput(gate.inputB)
			case XOR:
				w.state[gate.output] = w.getOutput(gate.inputA) ^ w.getOutput(gate.inputB)
			}
			return w.state[gate.output]
		}
	}
	return 0
}
