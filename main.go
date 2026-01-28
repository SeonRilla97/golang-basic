package main

import (
	"fmt"
	"time"
)

type Result struct {
	ID    int
	Value int
}

func worker(id int, done <-chan struct{}, jobs <-chan int) {
	for {
		select {
		case <-done:
			fmt.Printf("워커 %d 종료\n", id)
			return
		case job := <-jobs:
			fmt.Printf("워커 %d: 작업 %d 처리\n", id, job)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	funcs := make([]func(), 3)

	for i := 0; i < 3; i++ {
		funcs[i] = func() {
			fmt.Println(i)
		}
	}

	for _, f := range funcs {
		f()
	}
	// 출력: 3 3 3 (예상: 0 1 2)
}
