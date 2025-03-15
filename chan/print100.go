package main

import (
	"fmt"
	"sync"
)

func Print100(cg chan int, routingId int, wg *sync.WaitGroup) {
	fmt.Printf("start go routing %d\n", routingId)
	defer wg.Done()
	for {
		if numb, ok := <-cg; ok {
			fmt.Printf("go routing %d print %d\n", routingId, numb)
			if numb >= 100 {
				close(cg)
				return
			}
			numb += 1
			cg <- numb
		} else {
			return
		}
	}
}

func main() {
	fmt.Printf("start main\n")
	var cg = make(chan int, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go Print100(cg, 1, &wg)
	go Print100(cg, 2, &wg)
	cg <- 1
	wg.Wait()
}
