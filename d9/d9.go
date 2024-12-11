package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type filemeta struct {
	id              int
	arrayStartIndex int
	filesCount      int
	emptySpace      int
	emptySpaceIndex int // index of the first empty space in the data. -1 if no empty space
	checksum        int
	data            []int // only for debugging / visualization purposes
}

// this function reads the input from the txt file and splits the rules and prints
func readInput() []filemeta {
	// read from txt file
	file, err := os.Open("input9.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var input []filemeta

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// split line every 2 characters
		arrayStartIndex := 0
		for i := 0; i < len(line); i += 2 {
			end := i + 2
			if end > len(line) {
				end = len(line)
			}
			intBlocks := make([]int, len(line[i:end]))
			for i, v := range strings.Split(line[i:end], "") {
				intBlocks[i], _ = strconv.Atoi(v)
			}
			if len(intBlocks) == 1 {
				intBlocks = append(intBlocks, 0) // add 0 for empty space for the last input line
			}

			id := i / 2

			data := make([]int, 0)
			for i := 0; i < intBlocks[0]; i++ {
				data = append(data, id)
			}
			for i := 0; i < intBlocks[1]; i++ {
				data = append(data, -1)
			}
			arrayIndexOccupied := intBlocks[0] + intBlocks[1]
			checksum := id * sumfromNtoM(arrayStartIndex, arrayStartIndex+intBlocks[0])

			input = append(input, filemeta{id: id, arrayStartIndex: arrayStartIndex, filesCount: intBlocks[0], emptySpace: intBlocks[1], emptySpaceIndex: arrayStartIndex + intBlocks[0], checksum: checksum, data: data})

			arrayStartIndex += arrayIndexOccupied
		}

	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return input
}

func main() {
	startTime := time.Now()
	input := readInput()

	p1(input)
	p2(input)

	endTime := time.Since(startTime)
	fmt.Printf("Day9 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day9 execution took: 39 ms (39177 µs)
}

// sum of 1 to n
func sumOfN(n int) int {
	return n * (n + 1) / 2
}

// sum of n to m, given n < m
func sumfromNtoM(n, m int) int {
	return sumOfN(m-1) - sumOfN(n-1)
}

func p1(_input []filemeta) {
	input := make([]filemeta, len(_input))
	copy(input, _input)
	checksum := 0

	startIdPtr := 0
	indexPtr := 0
	endIdPtr := len(input) - 1
	for startIdPtr <= endIdPtr {
		// sum the startIdPtr
		filesCount := input[startIdPtr].filesCount
		checksum += (sumfromNtoM(indexPtr, indexPtr+filesCount) * startIdPtr)
		indexPtr += filesCount

		if startIdPtr == endIdPtr {
			break
		}

		// fill in the empty space
		emptySpaces := input[startIdPtr].emptySpace
		endFilesCount := input[endIdPtr].filesCount
		for i := 0; i < emptySpaces; i++ {
			if endFilesCount == 0 {
				if startIdPtr == endIdPtr-1 {
					break
				}
				endIdPtr--
				endFilesCount = input[endIdPtr].filesCount
			}
			checksum += (endIdPtr * indexPtr)
			indexPtr++
			endFilesCount--
			input[endIdPtr].filesCount = endFilesCount
		}
		startIdPtr++
	}

	fmt.Println("Checksum: ", checksum)
}

// used for debugging
func printFilemeta(input []filemeta) {
	for _, v := range input {
		fmt.Printf("%+v\n", v)
	}
	fmt.Println()
}

func p2(_input []filemeta) {
	input := make([]filemeta, len(_input))
	copy(input, _input)

	for i := len(input) - 1; i > 0; i-- {
		currentFileToMove := input[i]
		// find first free space to move the files
		freeSpaceIndex := findFirstFreeSpace(input, currentFileToMove)
		if freeSpaceIndex != -1 {
			// move the files
			moveFilesIntoFreeSpace(input, currentFileToMove, input[freeSpaceIndex])
		}
	}

	// sum the checksum
	checksum := 0
	for _, v := range input {
		checksum += v.checksum
	}
	fmt.Println("Checksum: ", checksum)

}

// find the first free space enough to fit the files from the LEFT to its current position
func findFirstFreeSpace(input []filemeta, currentFileToMove filemeta) int {
	maxSearchIndex := currentFileToMove.id
	for i, v := range input[:maxSearchIndex] {
		if v.emptySpace >= currentFileToMove.filesCount {
			return i
		}
	}
	return -1
}

// move the files into the free space
func moveFilesIntoFreeSpace(input []filemeta, filesToMove filemeta, freeSpace filemeta) {

	// update the data
	emptySpaceIndexInData := freeSpace.emptySpaceIndex - freeSpace.arrayStartIndex
	for i := emptySpaceIndexInData; i < emptySpaceIndexInData+filesToMove.filesCount; i++ {
		freeSpace.data[i] = filesToMove.id
	}

	// update the checksum of freeSpace since index has changed
	newChecksum := 0
	for i := 0; i < filesToMove.filesCount; i++ {
		newChecksum += filesToMove.id * (freeSpace.emptySpaceIndex + i)
	}

	// move the files
	freeSpace = filemeta{
		id:              freeSpace.id,
		arrayStartIndex: freeSpace.arrayStartIndex,
		filesCount:      freeSpace.filesCount,
		emptySpace:      freeSpace.emptySpace - filesToMove.filesCount,
		emptySpaceIndex: freeSpace.emptySpaceIndex + filesToMove.filesCount,
		checksum:        freeSpace.checksum + newChecksum,
		data:            freeSpace.data,
	}

	// update files to move
	for i := 0; i < filesToMove.filesCount; i++ {
		filesToMove.data[i] = -1 // empty space
	}
	// update the checksum of filesToMove since index has changed
	newChecksum2 := 0
	for i := 0; i < filesToMove.filesCount; i++ {
		newChecksum2 += filesToMove.id * (filesToMove.arrayStartIndex + i)
	}

	filesToMove = filemeta{
		id:              filesToMove.id,
		arrayStartIndex: filesToMove.arrayStartIndex,
		filesCount:      0,
		emptySpace:      filesToMove.filesCount + filesToMove.emptySpace,
		emptySpaceIndex: filesToMove.emptySpaceIndex - filesToMove.filesCount,
		checksum:        filesToMove.checksum - newChecksum2,
		data:            filesToMove.data,
	}

	input[freeSpace.id] = freeSpace
	input[filesToMove.id] = filesToMove
}
