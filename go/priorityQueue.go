// This example demonstrates a Priority queue built using the heap interface.
// Source https://golang.org/pkg/container/heap/
package main

import (
	"container/heap"
)

// An item is something we manage in a Priority queue.
type item struct {
	Value    path    // The Value of the item; arbitrary.
	Priority float64 // The Priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A priorityQueue implements heap.Interface and holds Items.
type priorityQueue []*item

func (pq priorityQueue) Len() int { return len(pq) }

func (pq priorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the Priority and Value of an item in the queue.
func (pq *priorityQueue) update(item *item, value path, priority float64) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.index)
}

// This example creates a priorityQueue with some items, adds and manipulates an item,
// and then removes the items in Priority order.
// func main() {
// 	// Some items and their priorities.
// 	items := map[string]float64{
// 		"banana": 3, "apple": 2, "pear": 4,
// 	}
//
// 	// Create a Priority queue, put the items in it, and
// 	// establish the Priority queue (heap) invariants.
// 	pq := make(priorityQueue, len(items))
// 	i := 0
// 	for value, priority := range items {
// 		pq[i] = &item{
// 			Value:    value,
// 			Priority: priority,
// 			index:    i,
// 		}
// 		i++
// 	}
// 	heap.Init(&pq)
//
// 	// Insert a new item and then modify its Priority.
// 	insert := &item{
// 		Value:    "orange",
// 		Priority: 1,
// 	}
// 	heap.Push(&pq, insert)
// 	pq.update(insert, insert.Value, 5)
//
// 	// Take the items out; they arrive in decreasing Priority order.
// 	for pq.Len() > 0 {
// 		item := heap.Pop(&pq).(*item)
// 		fmt.Printf("%.2f:%s ", item.Priority, item.Value)
// 	}
// }
