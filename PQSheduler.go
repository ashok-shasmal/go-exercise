package main

import (
	"container/heap"
	"fmt"
	"sync"
)

type PTask struct {
	Priority int
	index    int
	Msg      string
}

type PriorityQueue []*PTask

func (pq *PriorityQueue) CreatePriorityQueue() {
	for i := 100; i >= 1; i-- {
		s := fmt.Sprintf("Hello!!, I am task no. %d", i)
		t := PTask{Priority: i, Msg: s}
		*pq = append(*pq, &t)
	}
}

func (pq *PriorityQueue) Len() int { return len(*pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Pop() interface{} {

	orig := *pq
	n := len(orig)
	task := orig[n-1]
	orig[n-1] = nil
	*pq = orig[0 : n-1]
	return task
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(PTask)
	item.index = n
	*pq = append(*pq, &item)
}

// main function
func main() {

	var wg sync.WaitGroup
	ptaskQ := make(chan PTask)

	// go routine to pick the task and print according to the priority
	execute := func(ch <-chan PTask, wg *sync.WaitGroup) {
		defer wg.Done()
		for ts := range ch {
			fmt.Printf("Priority : %d, Message : %s\n", ts.Priority, ts.Msg)
		}
	}

	//Worker pool
	createWorkerpool := func(numOfWorkers int) {
		for i := 0; i < numOfWorkers; i++ {
			wg.Add(1)
			go execute(ptaskQ, &wg)
		}
	}

	createWorkerpool(4)

	//create the Priorty queue and fill items
	pq := make(PriorityQueue, 0)
	pq.CreatePriorityQueue()

	heap.Init(&pq)

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PTask)
		ptaskQ <- *item
	}

	close(ptaskQ)
	wg.Wait()

}
