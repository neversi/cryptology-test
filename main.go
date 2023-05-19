package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	in := make(chan []byte, 100)

	out, done := process(in)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(1 * time.Second)
		in <- ssmin
		time.Sleep(1 * time.Second)
		in <- []byte(update1)
		time.Sleep(1 * time.Second)
		in <- ss1
		time.Sleep(1 * time.Second)
		in <- []byte(update1)
	}()

	go func() {
		wg.Wait()
		fmt.Println("Shutdown process")
		close(in)
	}()

	for data := range out {
		fmt.Println(string(data))
	}

	<-done
	fmt.Println("Successfully shutdowned")
}
