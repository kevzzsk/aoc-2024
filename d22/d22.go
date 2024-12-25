package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func readInput() []int {
	// read from txt file
	file, err := os.Open("input22.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	result := []int{}

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, _ := strconv.Atoi(line)
		result = append(result, num)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

type bananaLedger struct {
	maxBanana   []bool
	totalBanana int
}

type secretBananaNumber struct {
	secretNumber int
	bananaPrice  int
	diff         int
}

func main() {
	startTime := time.Now()

	result := readInput()

	totalSum := 0
	secretNumberArr := make([][]secretBananaNumber, 0, len(result))
	for _, each := range result {
		secretNumber := each
		currArr := []secretBananaNumber{{secretNumber, getBananaPrice(secretNumber), 0}}
		for i := 0; i < 2000; i++ {
			nextSecretNumber := getNextSecretNumber(secretNumber)
			diff := getBananaPriceDiff(nextSecretNumber, secretNumber)
			currArr = append(currArr, secretBananaNumber{nextSecretNumber, getBananaPrice(nextSecretNumber), diff})
			secretNumber = nextSecretNumber
		}
		secretNumberArr = append(secretNumberArr, currArr)
		totalSum += secretNumber
	}

	fmt.Println("Total sum is: ", totalSum)

	blmap := make(map[string]bananaLedger)
	for buyerIndex, each := range secretNumberArr {
		for i := 3; i < len(each); i++ {
			currSecretBananaNumber := each[i]
			patternKey := fmt.Sprintf("%d,%d,%d,%d", each[i-3].diff, each[i-2].diff, each[i-1].diff, each[i].diff)
			if bl, ok := blmap[patternKey]; ok {
				// if found, update the existing one
				if !blmap[patternKey].maxBanana[buyerIndex] {
					bl.maxBanana[buyerIndex] = true
					bl.totalBanana += currSecretBananaNumber.bananaPrice
					// assign bl back to blmap
					blmap[patternKey] = bl
				}

			} else {
				// if not found, create a new one
				bl := bananaLedger{
					maxBanana:   make([]bool, len(secretNumberArr)),
					totalBanana: currSecretBananaNumber.bananaPrice,
				}
				bl.maxBanana[buyerIndex] = true
				blmap[patternKey] = bl
			}
		}
	}

	// get the max total banana
	maxTotalBanana := 0
	for _, each := range blmap {
		if each.totalBanana > maxTotalBanana {
			maxTotalBanana = each.totalBanana
		}
	}

	fmt.Println("Max total bananas: ", maxTotalBanana)

	endTime := time.Since(startTime)
	fmt.Printf("Day22 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day22 execution took: 1459 ms (1459203 µs)
}

func getBananaPrice(price int) int {
	return price % 10
}

func getBananaPriceDiff(newPrice, oldPrice int) int {
	return getBananaPrice(newPrice) - getBananaPrice(oldPrice)
}

func mixSecretNumbers(secretNumber, resultant int) int {
	return secretNumber ^ resultant
}

func pruneSecretNumber(secretNumber int) int {
	return secretNumber % 16777216
}

func getNextSecretNumber(secretNumber int) int {
	firsrOperation := pruneSecretNumber(mixSecretNumbers(secretNumber, secretNumber*64))
	secondOperation := pruneSecretNumber(mixSecretNumbers(firsrOperation, firsrOperation/32))
	thirdOperation := pruneSecretNumber(mixSecretNumbers(secondOperation, secondOperation*2048))
	return thirdOperation
}
