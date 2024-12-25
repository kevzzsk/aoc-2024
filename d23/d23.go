package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type connection struct {
	computerA string
	computerB string
}

func readInput() []connection {
	// read from txt file
	file, err := os.Open("input23.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var result []connection

	// read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		computers := strings.Split(line, "-")
		result = append(result, connection{computers[0], computers[1]})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

type node struct {
	computer      string
	nameContainsT bool
	connections   []*node
	connectionMap map[string]*node
}

type network struct {
	nodes []*node
}

type setOfThree struct {
	nodes []*node
}

func main() {
	startTime := time.Now()

	connections := readInput()

	// create a network of nodes
	network := network{}
	for _, connection := range connections {
		computerA := connection.computerA
		computerB := connection.computerB
		var nodeA, nodeB *node
		nodeA = network.getNodeByName(computerA)
		nodeB = network.getNodeByName(computerB)
		if nodeA == nil {
			nodeA = &node{computerA, strings.HasPrefix(computerA, "t"), []*node{}, make(map[string]*node)}
			network.nodes = append(network.nodes, nodeA)
		}
		if nodeB == nil {
			nodeB = &node{computerB, strings.HasPrefix(computerB, "t"), []*node{}, make(map[string]*node)}
			network.nodes = append(network.nodes, nodeB)
		}
		nodeA.connections = append(nodeA.connections, nodeB)
		nodeB.connections = append(nodeB.connections, nodeA)
		nodeA.connectionMap[computerB] = nodeB
		nodeB.connectionMap[computerA] = nodeA

	}

	setsOfThrees := make(map[string]bool)
	for _, n := range network.nodes {
		if n.nameContainsT {
			for _, connection := range n.connections {
				commonNodes := findAllCommonNeighbour(n, connection)
				for _, commonNode := range commonNodes {
					if commonNode != nil {
						currSetOfThree := &setOfThree{[]*node{n, connection, commonNode}}
						sort.Sort(currSetOfThree)
						setsOfThrees[currSetOfThree.String()] = true
					}
				}
			}
		}
	}
	fmt.Println("Number of sets of three: ", len(setsOfThrees))

	// part 2
	maxClique := make(map[string]*node)
	for _, n := range network.nodes {
		currClique := network.dfs(n, make(map[string]*node))
		if len(currClique) > len(maxClique) {
			maxClique = currClique
		}
	}

	// put all the nodes in maxClique into a set
	maxCliqueSet := &setOfThree{}
	for _, n := range maxClique {
		maxCliqueSet.nodes = append(maxCliqueSet.nodes, n)
	}
	sort.Sort(maxCliqueSet)

	fmt.Println("Max clique size: ", maxCliqueSet)

	endTime := time.Since(startTime)
	fmt.Printf("Day23 execution took: %v ms (%v µs)\n", endTime.Milliseconds(), endTime.Microseconds()) // Day23 execution took: 12 ms (12565 µs)
}

func (net *network) dfs(n *node, clique map[string]*node) map[string]*node {

	belongsInClique := true
	for _, cliqueNode := range clique {
		if n.connectionMap[cliqueNode.computer] == nil {
			belongsInClique = false
			break
		}
	}
	maxClique := clique
	if belongsInClique {
		// add node into clique
		clique[n.computer] = n
		maxClique = clique
		// try add his friends
		for _, connection := range n.connections {
			// only from those not already in the clique
			if _, ok := clique[connection.computer]; !ok {
				newclique := net.dfs(connection, clique)
				if len(newclique) > len(maxClique) {
					maxClique = newclique
				}
			}
		}
	}

	return maxClique
}

func (node *node) String() string {
	return node.computer
}

func (set *setOfThree) String() string {
	nodeNames := make([]string, 0)
	for _, n := range set.nodes {
		nodeNames = append(nodeNames, n.computer)
	}
	return strings.Join(nodeNames, ",")
}

func (set *setOfThree) Len() int {
	return len(set.nodes)
}

func (set *setOfThree) Less(i, j int) bool {
	return set.nodes[i].computer < set.nodes[j].computer
}

func (set *setOfThree) Swap(i, j int) {
	set.nodes[i], set.nodes[j] = set.nodes[j], set.nodes[i]
}

func (n *network) getNodeByName(computerName string) *node {
	for _, n := range n.nodes {
		if n.computer == computerName {
			return n
		}
	}
	return nil
}

func findAllCommonNeighbour(node1 *node, node2 *node) []*node {
	commonNeighbours := make([]*node, 0)
	for _, connection := range node1.connections {
		for _, connection2 := range node2.connections {
			if connection == connection2 {
				commonNeighbours = append(commonNeighbours, connection)
			}
		}
	}
	return commonNeighbours
}
