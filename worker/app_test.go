package worker

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	ch := make(chan int)
	go func() {
		ch <- 2
		ch <- 1
	}()
	go func() {
		res := <-ch
		fmt.Println(res)
		res2 := <-ch
		fmt.Println(res2)
		res3 := <-ch
		fmt.Println(res3)
	}()
	time.Sleep(2 * time.Second)
}

func TestNewWorker(t *testing.T) {
	w := NewWorker(2)
	list := make([]string, 0)
	ctx := context.Background()
	list = append(list, "1111", "2222", "3333")
	w.SetFunc(func(ctx context.Context, item string) error {
		fmt.Println(item)
		time.Sleep(2 * time.Second)
		return nil
	})
	for _, v := range list {
		w.Enqueue(v)
	}
	fmt.Println(w.Handle(ctx))
}
