package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	noBufChan := make(chan int)
	oneBufChan := make(chan int, 1)
	fmt.Printf("noBufChan: cap %d, len %d\n", cap(noBufChan), len(noBufChan))
	fmt.Printf("oneBufChan: cap %d, len %d\n", cap(oneBufChan), len(oneBufChan))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("start get noBufChan")
		a := <-noBufChan
		fmt.Println("end get noBufChan: a ", a)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("start get oneBufChan")
		a := <-oneBufChan
		fmt.Println("get oneBufChan: a ", a)
		b := <-oneBufChan
		fmt.Println("end get oneBufChan: b ", b)
	}()
	tick := time.Tick(1 * time.Second)
	<-tick
	fmt.Println("tick 1s")
	noBufChan <- 1

	<-tick
	fmt.Println("tick 2s")
	oneBufChan <- 11

	<-tick
	fmt.Println("tick 3s")
	oneBufChan <- 12

	wg.Wait()
}
