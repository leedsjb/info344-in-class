package main

import (
	"fmt"
	"math/rand"
	"time"
)

//someLongFunc is a function that might
//take a while to complete, so we want
//to run it on its own go routine
func someLongFunc(ch chan int) { //channel of ints called chan

	r := rand.Intn(2000) // create rand int between 0 and param
	d := time.Duration(r)
	time.Sleep(time.Millisecond * d)
	ch <- r // write to channel

}

func main() {
	//TODO:
	//create a channel and call
	//someLongFunc() on a go routine
	//passing the channel so that
	//someLongFunc() can communicate
	//its results

	rand.Seed(time.Now().UnixNano()) // convert time now to nanoseconds
	fmt.Println("Start long running channel...")

	n := 10
	ch := make(chan int, n) // make channel buffered with second param
	start := time.Now()
	for i := 0; i < n; i++ { // 10 iterations
		go someLongFunc(ch) // call someLongFunc n times, with x cores the computer run x functions in parallel
	}
	for i := 0; i < n; i++ {
		result := <-ch
		fmt.Printf("result was %d\n", result)
	}
	fmt.Printf("took %v\n", time.Since(start))
}
