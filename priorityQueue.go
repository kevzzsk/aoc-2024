package utils

import (
	"container/heap"
	"fmt"
)

// An Item is something we manage in a priority queue.
type Item[T any] struct {
	Value    T   // The value of the item; arbitrary.
	Priority int // The priority of the item in the queue.
	// The Index is needed by update and is maintained by the heap.Interface methods.
	Index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T any] []*Item[T]

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the lowest, priority so we use less than here
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue[T]) Push(x any) {
	n := len(*pq)
	item := x.(*Item[T])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue[T]) Update(item *Item[T], value T, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.Index)
}

func (pq *PriorityQueue[T]) Peek() *Item[T] {
	return (*pq)[0]
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return len(*pq) == 0
}

// print the priority queue in order
func (pq *PriorityQueue[T]) Print() {
	// Create a copy of the priority queue.
	copyPQ := make(PriorityQueue[T], pq.Len())
	copy(copyPQ, *pq)
	heap.Init(&copyPQ)

	// Pop elements from the copy and print them.
	for copyPQ.Len() > 0 {
		item := heap.Pop(&copyPQ).(*Item[T])
		fmt.Printf("%.2d:%v \n", item.Priority, item.Value)
	}
	fmt.Println()
}
